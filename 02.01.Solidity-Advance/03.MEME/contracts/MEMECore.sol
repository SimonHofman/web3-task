// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

import {MEMEToken} from "./MEMEToken.sol";
import {IBondingCurveParams} from "./interfaces/IBondingCurveParams.sol";
import {IMEMECore} from "./interfaces/IMEMECore.sol";
import {IMEMEFactory} from "./interfaces/IMEMEFactory.sol";
import {IMEMEHelper} from "./interfaces/IMEMEHelper.sol";
import {IMEMEVesting} from "./interfaces/IMEMEVesting.sol";
import {IVestingParams} from "./interfaces/IVestingParams.sol";
import {AccessControlUpgradeable} from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

contract MEMECore is IMEMECore, Initializable, UUPSUpgradeable, AccessControlUpgradeable, PausableUpgradeable, ReentrancyGuardUpgradeable {
    using  SafeERC20 for IERC20;
    using ECDSA for bytes32;

    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 public constant SIGNER_ROLE = keccak256("SIGNER_ROLE");
    bytes32 public constant DEPLOYER_ROLE = keccak256("DEPLOYER_ROLE");
    bytes32 public constant PAUSER_ROLE = keccak256("PAUSER_ROLE");

    uint256 public constant REQUEST_EXPIRY = 3600;
    uint256 public constant MIN_LIQUIDITY = 10 ether;
    uint256 public constant MAX_INITIAL_BUY_PERCENTAGE = 9990;

    uint256 public creationFee;
    uint256 public preBuyFeeRate;
    uint256 public tradingFeeRate;
    uint256 public graduationPlatformFeeRate;
    uint256 public graduationCreatorFeeRate;
    uint256 public minLockTime;

    IMEMEFactory public factory;
    IMEMEHelper public helper;
    IMEMEVesting public vesting;
    address  public platformFeeReceiver;
    address  public marginReceiver;
    uint256 public CHAIN_ID;

    mapping(address => TokenInfo) public tokenInfo;
    mapping(address => BondingCurveParams) public bondingCurve;
    mapping(bytes32 => bool) public usedRequestIds;
    address  public graduateFeeReceiver;

    modifier validToken(address token) {
        if (tokenInfo[token].creator == address(0)) revert InvalidCreatorParameters();
        _;
    }

    modifier onlyTradingToken(address token) {
        _onlyTradingToken(token);
        _;
    }

    constructor() {
        _disableInitializers();
    }

    function initialize(
        address _factory,
        address _helper,
        address _signer,
        address _platformFeeReceiver,
        address _marginReceiver,
        address _graduateFeeReceiver,
        address _admin
    ) public initializer {
        __AccessControl_init();
        __Pausable_init();
        __ReentrancyGuard_init();
        __UUPSUpgradeable_init();

        CHAIN_ID = block.chainid;
        factory = IMEMEFactory(_factory);
        helper = IMEMEHelper(_helper);
        platformFeeReceiver = _platformFeeReceiver;
        marginReceiver = _marginReceiver;
        graduateFeeReceiver = _graduateFeeReceiver;

        _grantRole(DEFAULT_ADMIN_ROLE, _admin);
        _grantRole(ADMIN_ROLE, _admin);
        _grantRole(SIGNER_ROLE, _signer);
        _grantRole(DEPLOYER_ROLE, _admin);
        _grantRole(PAUSER_ROLE, _admin);

        creationFee = 0.05 ether;
        preBuyFeeRate = 300;
        tradingFeeRate = 100;
        graduationPlatformFeeRate = 550;
        graduationCreatorFeeRate = 250;
        minLockTime = 86400;
    }

    function createToken(
        bytes calldata data,
        bytes calldata signature
    ) external payable nonReentrant whenNotPaused returns (address tokenAddress) {
        if (msg.value < creationFee) revert InsufficientFee();
        IMEMECore.CreateTokenParams memory params = abi.decode(data, (IMEMECore.CreateTokenParams));

        bytes32 messageHash = keccak256(abi.encodePacked(data, CHAIN_ID, address(this)));
        address signer = messageHash.recover(signature);
        if (!hasRole(SIGNER_ROLE, signer)) revert InvalidSigner();

        if (block.timestamp > params.timestamp + REQUEST_EXPIRY) revert RequestExpired();
        if (usedRequestIds[params.requestId]) revert RequestAlreadyProcessed();

        if (params.saleAmount > params.totalSupply) revert InvalidSaleParameters();
        if (params.saleAmount == 0) revert InvalidSaleParameters();
        if (params.saleAmount < params.totalSupply * params.initialBuyPercentage / 10000) revert InvalidSaleParameters();

        if (params.initialBuyPercentage > MAX_INITIAL_BUY_PERCENTAGE) revert InvalidInitialBuyPercentage();

        uint256 totalPaymentRequired = creationFee;
        uint256 initialTokens = 0;
        uint256 initialBNB = 0;
        uint256 adjustedBNBReserve = params.virtualBNBReserve;
        uint256 adjustedTokenReserve = params.virtualTokenReserve;
        uint256 preBuyFee;

        if (params.marginBnb > 0) {
            totalPaymentRequired += params.marginBnb;
        }

        if (params.initialBuyPercentage > 0) {
            (initialTokens, initialBNB, adjustedBNBReserve, adjustedTokenReserve) =
            _calculateInitialBuy(
                params.totalSupply,
                params.virtualBNBReserve,
                params.virtualTokenReserve,
                params.initialBuyPercentage
            );

            preBuyFee = (initialBNB * preBuyFeeRate) / 10000;
            totalPaymentRequired += initialBNB + preBuyFee;
        }

        if (msg.value < totalPaymentRequired) revert InsufficientFee();

        usedRequestIds[params.requestId] = true;

        tokenAddress = factory.deployToken(
            params.name,
            params.symbol,
            params.totalSupply,
            params.timestamp,
            params.nonce
        );

        address pair = helper.getPairAddress(tokenAddress);
        if (pair == address(0)) revert InvalidPair();
        MEMEToken(tokenAddress).setPair(pair);

        bondingCurve[tokenAddress] = BondingCurveParams({
            virtualBNBReserve: adjustedBNBReserve,
            virtualTokenReserve: adjustedTokenReserve,
            k: params.virtualBNBReserve * params.virtualTokenReserve,
            availableTokens: params.saleAmount - initialTokens,
            collectedBNB: initialBNB
        });

        tokenInfo[tokenAddress] = TokenInfo({
            creator: params.creator,
            createAt: block.timestamp,
            launchTime: params.launchTime,
            status: TokenStatus.TRADING,
            liquidityPool: pair
        });

        MEMEToken(tokenAddress).setTransferMode(MEMEToken.TransferMode.MODE_TRANSFER_CONTROLED);

        if (address(vesting) != address(0)) {
            MEMEToken(tokenAddress).setVestingContract(address(vesting));
        }

        if (initialTokens > 0) {
            if (params.vestingAllocations.length > 0 && address(vesting) != address(0)) {
                uint256 tokensToTransfer = _createVestingSchedules(
                    tokenAddress,
                    params.creator,
                    initialTokens,
                    params.initialBuyPercentage,
                    params.totalSupply,
                    params.vestingAllocations
                );

                if (tokensToTransfer > 0) {
                    IERC20(tokenAddress).safeTransfer(params.creator, tokensToTransfer);
                }
            } else {
                IERC20(tokenAddress).safeTransfer(params.creator, initialTokens);
            }

            emit TokenCreatedWithInitialBuy(
                tokenAddress,
                params.creator,
                initialTokens,
                initialBNB,
                params.initialBuyPercentage
            );
        }

        if (params.marginBnb > 0) {
            if (marginReceiver == address(0)) revert MarginReceiverNotSet();
            payable(marginReceiver).transfer(params.marginBnb);

            emit MarginDeposited(tokenAddress, params.creator, params.marginBnb, params.marginTime);
        }

        _sendValue(platformFeeReceiver, preBuyFee);
        _sendValue(platformFeeReceiver, creationFee);

        if (msg.value > totalPaymentRequired) {
            payable(msg.sender).transfer(msg.value - totalPaymentRequired);
        }

        emit TokenCreated(tokenAddress, params.creator, params.name, params.symbol, params.totalSupply, params.requestId);
    }

    function buy(
        address token,
        uint256 minTokenAmount,
        uint256 deadline
    ) external payable nonReentrant whenNotPaused validToken(token) onlyTradingToken(token) {
        if (block.timestamp > deadline || deadline >= block.timestamp + 1 days) revert TransactionExpired();
        if (msg.value == 0) revert InvalidNativeAmount();

        uint256 tradingFee = (msg.value * tradingFeeRate) / 10000;
        uint256 netBNBCount = msg.value - tradingFee;
        uint256 tokenAmount = helper.calculateTokenAmountOut(netBNBCount, bondingCurve[token]);

        if (tokenAmount > bondingCurve[token].availableTokens) {
            tokenAmount = bondingCurve[token].availableTokens;
            netBNBCount = helper.calculateRequiredBNB(tokenAmount, bondingCurve[token]);
            tradingFee = (netBNBCount * tradingFeeRate) / (10000 - tradingFeeRate);
            uint256 actualPayment = netBNBCount + tradingFee;

            if (msg.value > actualPayment) {
                payable(msg.sender).transfer(msg.value - actualPayment);
            }
        }
        if (tokenAmount < minTokenAmount) revert SlippageExceeded();

        bondingCurve[token].virtualBNBReserve += netBNBCount;
        bondingCurve[token].virtualTokenReserve -= tokenAmount;
        bondingCurve[token].availableTokens -= tokenAmount;
        bondingCurve[token].collectedBNB += netBNBCount;

        _sendValue(platformFeeReceiver, tradingFee);
        IERC20(token).safeTransfer(msg.sender, tokenAmount);

        if (bondingCurve[token].availableTokens < MIN_LIQUIDITY) {
            _changeTokenStatus(token, TokenStatus.PENDING_GRADUATION);
            MEMEToken(token).setTransferMode(MEMEToken.TransferMode.MODE_TRANSFER_RESTRICTED);
        }

        emit TokenBought(
            token,
            msg.sender,
            netBNBCount,
            tokenAmount,
            tradingFee,
            bondingCurve[token].virtualBNBReserve,
            bondingCurve[token].virtualTokenReserve,
            bondingCurve[token].availableTokens,
            bondingCurve[token].collectedBNB
        );
    }

    function sell(
        address token,
        uint256 tokenAmount,
        uint256 minBNBAmount,
        uint256 deadline
    ) external nonReentrant whenNotPaused validToken(token) onlyTradingToken(token) {
        if (block.timestamp > deadline) revert TransactionExpired();
        if (tokenAmount == 0) revert InvalidParameters();

        if (IERC20(token).balanceOf(msg.sender) < tokenAmount) revert InsufficientBalance();

        uint256 bnbAmount = helper.calculateBNBAmountOut(tokenAmount, bondingCurve[token]);
        uint256 tradingFee = (bnbAmount * tradingFeeRate) / 10000;
        uint256 netBNBAmount = bnbAmount - tradingFee;

        if (netBNBAmount < minBNBAmount) revert SlippageExceeded();
        if (bnbAmount > bondingCurve[token].collectedBNB) revert InsufficientBalance();

        IERC20(token).safeTransferFrom(msg.sender, address(this), tokenAmount);

        bondingCurve[token].virtualBNBReserve -= bnbAmount;
        bondingCurve[token].virtualTokenReserve += tokenAmount;
        bondingCurve[token].availableTokens += tokenAmount;
        bondingCurve[token].collectedBNB -= bnbAmount;

        _sendValue(platformFeeReceiver, tradingFee);
        _sendValue(msg.sender, netBNBAmount);

        emit TokenSold(
            token,
            msg.sender,
            tokenAmount,
            netBNBAmount,
            tradingFee,
            bondingCurve[token].virtualBNBReserve,
            bondingCurve[token].virtualTokenReserve,
            bondingCurve[token].availableTokens,
            bondingCurve[token].collectedBNB
        );
    }

    function graduateToken(address token) external onlyRole(DEPLOYER_ROLE) validToken(token) nonReentrant {
        TokenInfo storage info = tokenInfo[token];

        BondingCurveParams storage curve = bondingCurve[token];

        uint256 collectBNB = curve.collectedBNB;
        uint256 remainingTokens = curve.availableTokens;

        uint256 platformFee = (collectBNB * graduationPlatformFeeRate) / 10000;
        uint256 creatorFee = (collectBNB * graduationCreatorFeeRate) / 10000;
        uint256 liquidityBNB = collectBNB - platformFee - creatorFee;

        uint256 tokenPlatformFee = remainingTokens * graduationPlatformFeeRate / 10000;
        uint256 tokenCreatorFee = remainingTokens * graduationCreatorFeeRate / 10000;
        uint256 liquidityTokens = remainingTokens - tokenPlatformFee - tokenCreatorFee;

        MEMEToken(token).setTransferMode(MEMEToken.TransferMode.MODE_NORMAL);

        require(IERC20(token).balanceOf(address(this)) > remainingTokens, "Insufficient token for liquidity");
        IERC20(token).approve(address(helper), liquidityTokens);
        uint256 liquidityResult = helper.addLiquidityV2{value: liquidityBNB}(token, liquidityBNB, liquidityTokens);

        _changeTokenStatus(token, TokenStatus.GRADUATED);

        _sendValue(graduateFeeReceiver, platformFee);
        if (tokenPlatformFee > 0) IERC20(token).safeTransfer(graduateFeeReceiver, tokenPlatformFee);
        _sendValue(info.creator, creatorFee);
        if (tokenCreatorFee > 0) IERC20(token).safeTransfer(info.creator, tokenCreatorFee);

        emit TokenGraduated(token, liquidityBNB, liquidityTokens, liquidityResult);
    }

    function pauseToken(address token) external onlyRole(PAUSER_ROLE) validToken(token) {
        _changeTokenStatus(token, TokenStatus.PAUSED);
        emit TokenPaused(token);
    }

    function unpauseToken(address token) external onlyRole(PAUSER_ROLE) validToken(token) {
        if (tokenInfo[token].status != TokenStatus.PAUSED) revert InvalidPausedStatus();
        _changeTokenStatus(token, TokenStatus.TRADING);
        emit TokenUnpaused(token);
    }

    function blacklistToken(address token) external onlyRole(ADMIN_ROLE) validToken(token) {
        _changeTokenStatus(token, TokenStatus.BLACKLISTED);
        emit TokenBlacklisted(token);
    }

    function removeFromBlacklist(address token) external onlyRole(ADMIN_ROLE) validToken(token) {
        if (tokenInfo[token].status != TokenStatus.BLACKLISTED) revert InvalidBlackListedStatus();
        _changeTokenStatus(token, TokenStatus.TRADING);
        emit TokenRemovedFromBlacklist(token);
    }

    function getTokenInfo(address token) external view returns (TokenInfo memory) {
        return tokenInfo[token];
    }

    function getBondingCurve(address token) external view returns (BondingCurveParams memory) {
        return bondingCurve[token];
    }

    function calculateBuyAmount(address token, uint256 bnbAmount) external view returns (uint256 tokenAmount) {
        BondingCurveParams memory curve = bondingCurve[token];
        tokenAmount = helper.calculateTokenAmountOut(bnbAmount, curve);
        if (tokenAmount > curve.availableTokens) {
            tokenAmount = curve.availableTokens;
        }
    }

    function calculateBuyAmountWithFee(address token, uint256 bnbAmount) external view returns (uint256 tokenOut, uint256 netBNB, uint256 feeBNB) {
        BondingCurveParams memory curve = bondingCurve[token];
        (tokenOut, netBNB, feeBNB) = helper.calculateTokenAmountOutWithFee(bnbAmount, curve, tradingFeeRate);
        if (tokenOut > curve.availableTokens) {
            tokenOut = curve.availableTokens;
            netBNB = helper.calculateRequiredBNB(tokenOut, curve);
            feeBNB = (netBNB * tradingFeeRate) / (10000 - tradingFeeRate);
        }
    }

    function calculateSellReturn(address token, uint256 tokenAmount) external view returns (uint256) {
        return helper.calculateBNBAmountOut(tokenAmount, bondingCurve[token]);
    }

    function calculateSellReturnWithFee(address token, uint256 tokenAmount) external view returns (uint256 netBNB, uint256 feeBNB) {
        (netBNB, feeBNB) = helper.calculateBNBAmountOutWithFee(tokenAmount, bondingCurve[token], tradingFeeRate);
    }

    function calculateInitialBuyBNB(
        uint256 totalSupply,
        uint256 virtualBNBReserve,
        uint256 virtualTokenReserve,
        uint256 percentageBP
    ) external view returns (uint256 totalPayment, uint256 preBuyFee) {
        if (percentageBP == 0) return (0, 0);
        if (percentageBP > MAX_INITIAL_BUY_PERCENTAGE) revert InvalidParameters();

        (, uint256 bnbRequired,,) = _calculateInitialBuy(totalSupply, virtualBNBReserve, virtualTokenReserve, percentageBP);
        preBuyFee = (bnbRequired * preBuyFeeRate) / 10000;
        totalPayment = bnbRequired + preBuyFee;
    }

    function _sendValue(address to, uint256 amount) internal {
        if (amount == 0) return;
        if (to == address(0)) {
            to = platformFeeReceiver;
            require(to != address(0), "Platform fee receiver not set");
        }
        uint32 size;
        assembly {
            size := extcodesize(to)
        }
        if (size > 0) {
            (bool ok,) = payable(to).call{value: amount}("");
            if (!ok) {
                (bool fallbackSuccess,) = payable(platformFeeReceiver).call{value: amount}("");
                require(fallbackSuccess, "BNB_SEND_FAILED_TO_FALLBACK");
                emit CreatorFeeRedirected(to, platformFeeReceiver, amount);
            }
        } else {
            (bool ok,) = payable(to).call{value: amount}("");
            require(ok, "BNB_SEND_FAILED");
        }
    }

    function _createVestingSchedules(
        address tokenAddress,
        address beneficiary,
        uint256 initialTokens,
        uint256 initialBuyPercentage,
        uint256 totalSupply,
        VestingAllocation[] memory vestingAllocations
    ) internal returns (uint256 tokensToTransfer) {
        uint256 totalVestedAmount;
        uint256 totalBurnedAmount;

        for (uint256 i = 0; i < vestingAllocations.length; i++) {
            if (vestingAllocations[i].amount == 0) {
                revert InvalidAmountParameters();
            }
            if (vestingAllocations[i].mode == VestingMode.BURN) {
                totalBurnedAmount += vestingAllocations[i].amount;
            } else {
                totalVestedAmount += vestingAllocations[i].amount;
                if (vestingAllocations[i].mode == VestingMode.LINEAR) {
                    if (vestingAllocations[i].duration == 0) {
                        revert InvalidDurationParameters();
                    }
                    if (vestingAllocations[i].duration < minLockTime) {
                        revert InvalidDurationParameters();
                    }
                }
            }
        }

        if (initialBuyPercentage < totalVestedAmount + totalBurnedAmount) revert InvalidVestingParameters();
        uint256 tokensToVest = (totalSupply + totalVestedAmount) / 10000;
        uint256 tokensToBurn = (totalSupply + totalBurnedAmount) / 10000;
        if (initialTokens < tokensToVest + tokensToBurn) revert InvalidParameters();

        tokensToTransfer = initialTokens - tokensToVest - tokensToBurn;
        if (tokensToBurn > 0) {
            MEMEToken(tokenAddress).burn(tokensToBurn);
            emit MEMETokensBurned(beneficiary, tokensToBurn);
        }

        if (tokensToVest > 0) {
            VestingAllocation[] memory actualVestingAllocation = new VestingAllocation[](vestingAllocations.length);

            uint256 allocatedTokens;
            TokenInfo memory info = tokenInfo[tokenAddress];
            uint256 actualLaunchTime = info.launchTime;
            if (actualLaunchTime == 0) {
                actualLaunchTime = block.timestamp;
            }
            int256 lastNonBurnIndex = - 1;
            for (uint256 i = 0; i < vestingAllocations.length; i++) {
                if (vestingAllocations[i].mode != VestingMode.BURN) {
                    lastNonBurnIndex = int256(i);
                }
            }
            for (uint256 i = 0; i < vestingAllocations.length; i++) {
                uint256 allocationAmount;
                if (vestingAllocations[i].mode == VestingMode.BURN) {
                    allocationAmount = 0;
                } else if (int256(i) == lastNonBurnIndex) {
                    allocationAmount = tokensToVest - allocatedTokens;
                } else {
                    allocationAmount = (totalSupply * vestingAllocations[i].amount) / 10000;
                    allocatedTokens += allocationAmount;
                }
                actualVestingAllocation[i] = VestingAllocation({
                    amount: allocationAmount,
                    lauchTime: actualLaunchTime,
                    duration: vestingAllocations[i].duration,
                    mode: vestingAllocations[i].mode
                });
            }

            IERC20(tokenAddress).approve(address(vesting), tokensToVest);

            vesting.createVestingSchedules(
                tokenAddress,
                beneficiary,
                actualVestingAllocation
            );

            emit VestingCreated(tokenAddress, beneficiary, tokensToVest, actualVestingAllocation.length);
        }

        return tokensToTransfer;
    }

    function _calculateInitialBuy(
        uint256 totalSupply,
        uint256 virtualBNBReserve,
        uint256 virtualTokenReserve,
        uint256 percentageBP
    ) internal pure returns (uint256 tokenOut, uint256 bnbRequired, uint256 newBNBReserve, uint256 newTokenReserve) {
        if (percentageBP > MAX_INITIAL_BUY_PERCENTAGE) revert InvalidPercentageBP();
        tokenOut = (totalSupply * percentageBP) / 10000;

        uint256 k = virtualBNBReserve * virtualTokenReserve;
        newTokenReserve = virtualTokenReserve - tokenOut;
        newBNBReserve = k / newTokenReserve;
        bnbRequired = newBNBReserve - virtualBNBReserve;

        return (tokenOut, bnbRequired, newBNBReserve, newTokenReserve);
    }

    function _changeTokenStatus(address token, TokenStatus newStatus) internal {
        TokenStatus oldStatus = tokenInfo[token].status;
        tokenInfo[token].status = newStatus;
        emit TokenStatusChanged(token, oldStatus, newStatus);
    }

    function _onlyTradingToken(address token) internal view {
        TokenInfo memory info = tokenInfo[token];
        if (info.status != TokenStatus.TRADING) {
            revert TokenNotTrading();
        }
        if (block.timestamp < info.launchTime) {
            revert TokenNotLauchedYet();
        }
    }

    function setPlatfromFeeReceiver(address _receiver) external onlyRole(ADMIN_ROLE) {
        if (_receiver == address(0)) revert ZeroAddress();
        address oldPlatformFeeReceiver = platformFeeReceiver;
        platformFeeReceiver = _receiver;
        emit PlatformFeeReceiverChanged(oldPlatformFeeReceiver, _receiver);
    }

    function setGraduateFeeReceiver(address _receiver) external onlyRole(ADMIN_ROLE) {
        if (_receiver == address(0)) revert ZeroAddress();
        address oldGraduateFeeReceiver = graduateFeeReceiver;
        graduateFeeReceiver = _receiver;
        emit GraduateFeeReceiverChanged(oldGraduateFeeReceiver, _receiver);
    }

    function setFactory(address _factory) external onlyRole(ADMIN_ROLE) {
        if (_factory == address(0)) revert ZeroAddress();
        address oldFactory = address(factory);
        factory = IMEMEFactory(_factory);
        emit FactoryChanged(oldFactory, _factory);
    }

    function setHelper(address _helper) external onlyRole(ADMIN_ROLE) {
        if (_helper == address(0)) revert ZeroAddress();
        address oldHelper = address(helper);
        helper = IMEMEHelper(_helper);
        emit HelperChanged(oldHelper, _helper);
    }

    function setVesting(address _vesting) external onlyRole(ADMIN_ROLE) {
        if (_vesting == address(0)) revert ZeroAddress();
        address oldVesting = address(vesting);
        vesting = IMEMEVesting(_vesting);
        emit VestingChanged(oldVesting, _vesting);
    }

    function setMarginReceiver(address _marginReceiver) external onlyRole(ADMIN_ROLE) {
        if (_marginReceiver == address(0)) revert ZeroAddress();
        address oldMarginReceiver = marginReceiver;
        marginReceiver = _marginReceiver;
        emit MarginReceiverChanged(oldMarginReceiver, _marginReceiver);
    }

    function setCreationFee(uint256 _fee) external onlyRole(ADMIN_ROLE) {
        if (_fee > 0.1 ether) revert InvalidAmountParameters();
        creationFee = _fee;
        emit CreationFeeChanged(_fee);
    }

    function setPreBuyFeeRate(uint256 _rate) external onlyRole(ADMIN_ROLE) {
        if (_rate > 600) revert InvalidAmountParameters();
        preBuyFeeRate = _rate;
        emit PreBuyFeeRateChanged(_rate);
    }

    function setTradingFeeRate(uint256 _rate) external onlyRole(ADMIN_ROLE) {
        if (_rate > 200) revert InvalidAmountParameters();
        tradingFeeRate = _rate;
        emit TradingFeeRateChanged(_rate);
    }

    function setGraduationFeeRates(uint256 _platformRate, uint256 _creatorRate) external onlyRole(ADMIN_ROLE) {
        if (_platformRate > 1100 || _creatorRate > 500) revert InvalidAmountParameters();
        graduationPlatformFeeRate = _platformRate;
        graduationCreatorFeeRate = _creatorRate;
        emit GraduationFeeRatesChanged(_platformRate, _creatorRate);
    }

    function setMinLockTime(uint256 _time) external onlyRole(ADMIN_ROLE) {
        minLockTime = _time;
        emit MinLockTimeChanged(_time);
    }

    function pause() external onlyRole(PAUSER_ROLE) {
        _pause();
    }

    function uppause() external onlyRole(PAUSER_ROLE) {
        _unpause();
    }

    function emergencyWithdraw(address token, uint256 amount) external onlyRole(ADMIN_ROLE) {
        if (token == address(0)) {
            payable(msg.sender).transfer(amount);
        } else {
            IERC20(token).safeTransfer(msg.sender, amount);
        }
    }

    function _authorizeUpgrade(address newImplementation) internal override onlyRole(ADMIN_ROLE) {}

    receive() external payable {}
}