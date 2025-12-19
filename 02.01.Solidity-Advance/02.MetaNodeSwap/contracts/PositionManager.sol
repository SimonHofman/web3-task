// SPDX-License-Identifier: GPL-2.0-or-later
pragma solidity ^0.8.24;
pragma abicoder v2;

import {IPool} from "./interfaces/IPool.sol";
import {IPoolManager} from "./interfaces/IPoolManager.sol";
import {IPositionManager} from "./interfaces/IPositionManager.sol";
import {FixedPoint128} from "./libraries/FixedPoint128.sol";
import {FullMath} from "./libraries/FullMath.sol";
import {LiquidityAmounts} from "./libraries/LiquidityAmounts.sol";
import {TickMath} from "./libraries/TickMath.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {ERC721} from "@openzeppelin/contracts/token/ERC721/ERC721.sol";

contract PositionManager is IPositionManager, ERC721 {
    IPoolManager public poolManager;

    uint176 private _nextId = 1;

    constructor(address _pooolManager) ERC721("MetaNodeSwapPosition", "MNSP") {
        poolManager = IPoolManager(_pooolManager);
    }

    mapping(uint256 => PositionInfo) public positions;

    function getAllPositions() external view override returns (PositionInfo[] memory positionInfo) {
        positionInfo = new PositionInfo[](_nextId - 1);
        for (uint32 i = 0; i < _nextId - 1; i++) {
            positionInfo[i] = positions[i + 1];
        }
        return positionInfo;
    }

    function getSender() public view returns (address) {
        return msg.sender;
    }

    function _blockTimestamp() internal view virtual returns (uint256) {
        return block.timestamp;
    }

    modifier checkDeadline(uint256 deadline) {
        require(deadline >= _blockTimestamp(), "Transaction too old");
        _;
    }

    function mint(
        MintParams calldata params
    ) external payable override checkDeadline(params.deadline) returns (
        uint256 positionId,
        uint128 liquidity,
        uint256 amount0,
        uint256 amount1
    ) {
        address _pool = poolManager.getPool(
            params.token0,
            params.token1,
            params.index
        );
        IPool pool = IPool(_pool);

        uint160 sqrtPriceX96 = pool.sqrtPriceX96();
        uint160 sqrtRatioAX96 = TickMath.getSqrtPriceAtTick(pool.tickLower());
        uint160 sqrtRatioBX96 = TickMath.getSqrtPriceAtTick(pool.tickUpper());

        liquidity = LiquidityAmounts.getLiquidityForAmounts(
            sqrtPriceX96,
            sqrtRatioAX96,
            sqrtRatioBX96,
            params.amount0Desired,
            params.amount1Desired
        );

        bytes memory data = abi.encode(params.token0, params.token1, params.index, msg.sender);

        (amount0, amount1) = pool.mint(address(this), liquidity, data);

        _mint(params.recipient, (positionId = _nextId++));

        (, uint256 feeGrowthInside0LastX128, uint256 feeGrowthInside1LastX128,,) = pool.getPosition(address(this));

        positions[positionId] = PositionInfo({
            id: positionId,
            owner: params.recipient,
            token0: params.token0,
            token1: params.token1,
            index: params.index,
            fee: pool.fee(),
            liquidity: liquidity,
            tickLower: pool.tickLower(),
            tickUpper: pool.tickUpper(),
            tokensOwed0: 0,
            tokensOwed1: 0,
            feeGrowthInside0LastX128: feeGrowthInside0LastX128,
            feeGrowthInside1LastX128: feeGrowthInside1LastX128
        });
    }

    modifier  isAuthorizedForToken(uint256 tokenId) {
        address owner = ERC721.ownerOf(tokenId);
        require(_isAuthorized(owner, msg.sender, tokenId), "Not approved");
        _;
    }

    function burn(uint256 tokenId) external override isAuthorizedForToken(tokenId) returns (uint256 amount0, uint256 amount1) {
        PositionInfo storage position = positions[tokenId];
        uint128 _liquidity = position.liquidity;
        address _pool = poolManager.getPool(position.token0, position.token1, position.index);
        IPool pool = IPool(_pool);
        (amount0, amount1) = pool.burn(_liquidity);
        (, uint256 feeGrowthInside0LastX128, uint256 feeGrowthInside1LastX128,,) = pool.getPosition(address(this));

        position.tokensOwed0 +=
            uint128(amount0) +
            uint128(
                FullMath.mulDiv(
                    feeGrowthInside0LastX128 - position.feeGrowthInside0LastX128,
                    position.liquidity,
                    FixedPoint128.Q128
                )
            );

        position.tokensOwed1 +=
            uint128(amount1) +
            uint128(
                FullMath.mulDiv(
                    feeGrowthInside1LastX128 - position.feeGrowthInside1LastX128,
                    position.liquidity,
                    FixedPoint128.Q128
                )
            );

        position.feeGrowthInside0LastX128 = feeGrowthInside0LastX128;
        position.feeGrowthInside1LastX128 = feeGrowthInside1LastX128;
        position.liquidity = 0;
    }

    function collect(
        uint256 positionId,
        address recipient
    ) external override isAuthorizedForToken(positionId) returns (uint256 amount0, uint256 amount1) {
        PositionInfo storage position = positions[positionId];
        address _pool = poolManager.getPool(position.token0, position.token1, position.index);
        IPool pool = IPool(_pool);
        (amount0, amount1) = pool.collect(recipient, position.tokensOwed0, position.tokensOwed1);

        position.tokensOwed0 = 0;
        position.tokensOwed1 = 0;

        if (position.liquidity == 0) {
            _burn(positionId);
        }
    }

    function mintCallback(uint256 amount0, uint256 amount1, bytes calldata data) external override {
        (address token0, address token1, uint32 index, address payer) = abi
            .decode(data, (address, address, uint32, address));
        address _pool = poolManager.getPool(token0, token1, index);
        require(_pool == msg.sender, "Invalid callback caller");

        if (amount0 > 0) {
            IERC20(token0).transferFrom(payer, msg.sender, amount0);
        }
        if (amount1 > 0) {
            IERC20(token1).transferFrom(payer, msg.sender, amount1);
        }
    }
}
