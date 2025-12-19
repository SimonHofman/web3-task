// SPDX-License-Identifier: GPL-2.0-or-later
pragma solidity ^0.8.24;

import {IPool} from "./interfaces/IPool.sol";
import {IPoolManager} from "./interfaces/IPoolManager.sol";
import {ISwapRouter} from "./interfaces/ISwapRouter.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

pragma abicoder v2;


contract SwapRouter is ISwapRouter {
    IPoolManager public poolManager;

    constructor(address _poolManager){
        poolManager = IPoolManager(_poolManager);
    }

    function parseRevertReason(
        bytes memory reason
    ) private pure returns (int256, int256) {
        if (reason.length != 64) {
            if (reason.length < 68) revert("Unexpected error");
            assembly{
                reason := add(reason, 0x04)
            }
            revert(abi.decode(reason, (string)));
        }
        return abi.decode(reason, (int256, int256));
    }

    function swapInPool(
        IPool pool,
        address recipient,
        bool zeroForOne,
        int256 amountSpecified,
        uint160 sqrtPriceLimitX96,
        bytes calldata data
    ) external returns (int256 amount0, int256 amount1) {
        try
        pool.swap(recipient, zeroForOne, amountSpecified, sqrtPriceLimitX96, data)
        returns (int256 _amoutn0, int256 _amount1){
            return (amount0, amount1);
        } catch Error(string memory reason) {
            return parseRevertReason(bytes(reason));
        }
    }

    function exactInput(
        ExactInputParams calldata params
    ) external payable override returns (uint256 amountOut) {
        uint256 amountIn = params.amountIn;
        bool zeroForOne = params.tokenIn < params.tokenOut;

        for (uint256 i = 0; i < params.indexPath.length; i++) {
            address poolAddress = poolManager.getPool(params.tokenIn, params.tokenOut, params.indexPath[i]);

            require(poolAddress != address(0), "Pool does not exist");

            IPool pool = IPool(poolAddress);

            bytes memory data = abi.encode(
                params.tokenIn,
                params.tokenOut,
                params.indexPath[i],
                params.recipient == address(0) ? address(0) : msg.sender
            );

            (int256 amount0, int256 amount1) = this.swapInPool(
                pool,
                params.recipient,
                zeroForOne,
                int256(amountIn),
                params.sqrtPriceLimitX96,
                data
            );

            amountIn -= uint256(zeroForOne ? amount0 : amount1);
            amountOut += uint256(zeroForOne ? - amount1 : - amount0);

            if (amountIn == 0) {
                break;
            }
        }

        require(amountOut >= params.amountOutMininum, "Slippage exceeded");

        emit Swap(msg.sender, zeroForOne, params.amountIn, amountIn, amountOut);

        return amountOut;
    }

    function exactOutput(
        ExactOutputParams calldata params
    ) external payable override returns (uint256 amountIn) {
        uint256 amountOut = params.amountOut;

        bool zeroForOne = params.tokenIn < params.tokenOut;

        for (uint256 i = 0; i < params.indexPath.length; i++) {
            address poolAddress = poolManager.getPool(params.tokenIn, params.tokenOut, params.indexPath[i]);

            require(poolAddress != address(0), "Pool not found");

            IPool pool = IPool(poolAddress);

            bytes memory data = abi.encode(
                params.tokenIn,
                params.tokenOut,
                params.indexPath[i],
                params.recipient == address(0) ? address(0) : msg.sender
            );

            (int256 amount0, int256 amount1) = this.swapInPool(
                pool,
                params.recipient,
                zeroForOne,
                - int256(amountOut),
                params.sqrtPriceLimitX96,
                data
            );

            amountOut -= uint256(zeroForOne ? - amount1 : - amount0);
            amountIn += uint256(zeroForOne ? amount0 : amount1);

            if (amountOut == 0) {
                break;
            }
        }

        require(amountIn <= params.amountInMaximum, "Slippage exceeded");

        emit Swap(msg.sender, zeroForOne, params.amountOut, amountOut, amountIn);

        return amountIn;
    }

    function quoteExactInput(
        QuoteExactInputParams calldata params
    ) external override returns (uint256 amountOut){
        return this.exactInput(
            ExactInputParams({
                tokenIn: params.tokenIn,
                tokenOut: params.tokenOut,
                indexPath: params.indexPath,
                recipient: address(0),
                deadline: block.timestamp + 1 hours,
                amountIn: params.amountIn,
                amountOutMininum: 0,
                sqrtPriceLimitX96: params.sqrtPriceLimitX96
            })
        );
    }

    function quoteExactOutput(
        QuoteExactOutputParams calldata params
    ) external override returns (uint256 amountIn){
        return this.exactOutput(
            ExactOutputParams({
                tokenIn: params.tokenIn,
                tokenOut: params.tokenOut,
                indexPath: params.indexPath,
                recipient: address(0),
                deadline: block.timestamp + 1 hours,
                amountOut: params.amountOut,
                amountInMaximum: type(uint256).max,
                sqrtPriceLimitX96: params.sqrtPriceLimitX96
            })
        );
    }

    function swapCallback(
        int256 amount0Delta,
        int256 amount1Delta,
        bytes calldata data
    ) external override {
        (address  tokenIn, address  tokenOut, uint256 index, address  payer) = abi
            .decode(data, (address, address, uint256, address));
        address _pool = poolManager.getPool(tokenIn, tokenOut, uint32(index));

        require(_pool == msg.sender, "Invalid callback caller");

        uint256 amountToPay = amount0Delta > 0 ? uint256(amount0Delta) : uint256(amount1Delta);

        if (payer == address(0)) {
            assembly {
                let ptr := mload(0x40)
                mstore(ptr, amount0Delta)
                mstore(add(ptr, 0x20), amount1Delta)
                revert(ptr, 64)
            }
        }

        if (amountToPay > 0) {
            IERC20(tokenIn).transferFrom(payer, _pool, amountToPay);
        }
    }
}
