// SPDX-License-Identifier: MIT
pragma solidity ^0.8.22;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {AccessControlUpgradeable} from "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import {PausableUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/PausableUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";

import {IMEMEVesting} from "./interfaces/IMEMEVesting.sol";
import {IVestingParams} from "./interfaces/IVestingParams.sol";

contract MEMEVesting is IMEMEVesting, Initializable, UUPSUpgradeable, AccessControlUpgradeable, PausableUpgradeable, ReentrancyGuardUpgradeable {
    using  SafeERC20 for IERC20;

    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 public constant OPERATOR_ROLE = keccak256("OPERATOR_ROLE");

    mapping(address => mapping(address => mapping(uint256 => VestingSchedule))) public vestingSchedules;

    mapping(address => mapping(address => uint256)) public scheduleCount;

    mapping(address => mapping(address => uint256)) public totalVestedAmount;

    mapping(address => uint256) public totalTokenLocked;

    constructor(){
        _disableInitializers();
    }

    function initialize(
        address _admin,
        address _operator
    ) public initializer {
        __AccessControl_init();
        __Pausable_init();
        __ReentrancyGuard_init();
        __UUPSUpgradeable_init();

        _grantRole(DEFAULT_ADMIN_ROLE, _admin);
        _grantRole(ADMIN_ROLE, _admin);
        _grantRole(OPERATOR_ROLE, _operator);
    }

    function createVestingSchedules(
        address token,
        address beneficiary,
        VestingAllocation[] calldata allocations
    ) external override onlyRole(OPERATOR_ROLE) whenNotPaused returns (uint256[] memory scheduleIds){
        if (token == address(0) || beneficiary == address(0)) revert InvalidAddressParameters();
        if (allocations.length == 0) revert InvalidLengthParameters();

        scheduleIds = new uint256[](allocations.length);
        uint256 totalAmount = 0;
        for (uint256 i = 0; i < allocations.length; i++) {
            if (allocations[i].mode != VestingMode.BURN) {
                if (allocations[i].amount == 0) revert InvalidAmountParameters();
                totalAmount += allocations[i].amount;
            }
        }
        if (totalAmount > 0) {
            IERC20(token).safeTransferFrom(msg.sender, address(this), totalAmount);
        }
        for (uint256 i = 0; i < allocations.length; i++) {
            uint256 scheduleId = scheduleCount[token][beneficiary];
            scheduleCount[token][beneficiary]++;

            uint256 startTime = allocations[i].lauchTime;
            if (startTime == 0) {
                startTime = block.timestamp;
            }
            uint256 endTime = startTime + allocations[i].duration;

            vestingSchedules[token][beneficiary][scheduleId] = VestingSchedule({
                totalAmount: allocations[i].amount,
                startTime: startTime,
                endTime: endTime,
                claimedAmount: 0,
                revoked: false,
                mode: allocations[i].mode
            });

            scheduleIds[i] = scheduleId;
            totalVestedAmount[token][beneficiary] += allocations[i].amount;
            if (allocations[i].mode != VestingMode.BURN) {
                totalTokenLocked[token] += allocations[i].amount;
            }

            emit VestingScheduleCreated(
                token,
                beneficiary,
                scheduleId,
                allocations[i].amount,
                startTime,
                endTime
            );
        }
    }

    function claim(
        address token,
        uint256 scheduleId
    ) public override nonReentrant whenNotPaused returns (uint256 claimableAmount) {
        VestingSchedule storage schedule = vestingSchedules[token][msg.sender][scheduleId];
        if (schedule.totalAmount == 0) revert ScheduleNotFound();
        if (schedule.revoked) revert ScheduleRevoked();
        claimableAmount = _calculateClaimableAmount(schedule);
        if (claimableAmount == 0) revert NoClaimableAmount();
        schedule.claimedAmount += claimableAmount;
        totalTokenLocked[token] -= claimableAmount;
        IERC20(token).safeTransfer(msg.sender, claimableAmount);
        emit TokensClaimed(token, msg.sender, scheduleId, claimableAmount);
    }

    function claimAll(
        address token
    ) external override nonReentrant whenNotPaused returns (uint256 totalClaimed) {
        uint256 count = scheduleCount[token][msg.sender];

        for (uint256 i = 0; i < count; i++) {
            VestingSchedule storage schedule = vestingSchedules[token][msg.sender][i];
            if (schedule.totalAmount == 0 || schedule.revoked) continue;
            uint256 claimableAmount = _calculateClaimableAmount(schedule);
            if (claimableAmount > 0) {
                schedule.claimedAmount += claimableAmount;
                totalTokenLocked[token] -= claimableAmount;
                totalClaimed += claimableAmount;
                emit TokensClaimed(token, msg.sender, i, claimableAmount);
            }
        }

        if (totalClaimed > 0) {
            IERC20(token).safeTransfer(msg.sender, totalClaimed);
        }
    }

    function _calculateClaimableAmount(
        VestingSchedule memory schedule
    ) private view returns (uint256) {
        if (schedule.mode == VestingMode.BURN) {
            return 0;
        }
        if (block.timestamp <= schedule.startTime) {
            return 0;
        }
        uint256 vestedAmount;
        if (block.timestamp >= schedule.endTime) {
            vestedAmount = schedule.totalAmount;
        } else {
            if (schedule.mode == VestingMode.ClIFF) {
                vestedAmount = 0;
            } else if (schedule.mode == VestingMode.LINEAR) {
                uint256 timePassed = block.timestamp - schedule.startTime;
                uint256 totalDuration = schedule.endTime - schedule.startTime;
                vestedAmount = (schedule.totalAmount * timePassed) / totalDuration;
            }
        }
        return vestedAmount > schedule.claimedAmount ? vestedAmount - schedule.claimedAmount : 0;
    }

    function getVestingSchedule(
        address token,
        address beneficiary,
        uint256 scheduleId
    ) external view override returns (VestingSchedule memory) {
        return vestingSchedules[token][beneficiary][scheduleId];
    }

    function getVestingScheduleCount(
        address token,
        address beneficiary
    ) external view override returns (uint256) {
        return scheduleCount[token][beneficiary];
    }

    function getClaimableAmount(
        address token,
        address beneficiary,
        uint256 scheduleId
    ) external view override returns (uint256) {
        VestingSchedule memory schedule = vestingSchedules[token][beneficiary][scheduleId];
        if (schedule.totalAmount == 0 || schedule.revoked) return 0;
        return _calculateClaimableAmount(schedule);
    }

    function getTotalClaimableAmount(
        address token,
        address beneficiary
    ) external view override returns (uint256 totalClaimable) {
        uint256 count = scheduleCount[token][beneficiary];
         for (uint256 i = 0; i < count; i++) {
            VestingSchedule memory schedule = vestingSchedules[token][beneficiary][i];
            if (schedule.totalAmount == 0 || schedule.revoked) continue;
            totalClaimable += _calculateClaimableAmount(schedule);
        }
    }

    function getTotalVestedAmount(
        address token,
        address beneficiary
    ) external view override returns (uint256 vested, uint256 claimed, uint256 locked) {
        uint256 count = scheduleCount[token][beneficiary];

        for (uint256 i = 0; i < count; i++) {
            VestingSchedule memory schedule = vestingSchedules[token][beneficiary][i];
            if (schedule.totalAmount > 0 && !schedule.revoked) {
                if (schedule.mode == VestingMode.BURN) {
                    continue;
                }
                uint256 vestedForSchedule;

                if (block.timestamp >= schedule.endTime) {
                    vestedForSchedule = schedule.totalAmount;
                } else if (block.timestamp > schedule.startTime) {
                    uint256 timePassed = block.timestamp - schedule.startTime;
                    uint256 totalDuration = schedule.endTime - schedule.startTime;
                    vestedForSchedule = schedule.totalAmount * timePassed / totalDuration;
                }

                vested += vestedForSchedule;
                claimed += schedule.claimedAmount;
                locked += (schedule.totalAmount - vestedForSchedule);
            }
        }

        return (vested, claimed, locked);
    }

    function revokeVestingSchedule(
        address token,
        address beneficiary,
        uint256 scheduleId
    ) external override onlyRole(ADMIN_ROLE) {
        VestingSchedule storage schedule = vestingSchedules[token][beneficiary][scheduleId];

        if (schedule.totalAmount == 0) revert ScheduleNotFound();
        if (schedule.revoked) revert ScheduleRevoked();
        uint256 remainingAmount;
        if (schedule.mode != VestingMode.BURN) {
            uint256 claimableAmount = _calculateClaimableAmount(schedule);
            if (claimableAmount > 0) {
                schedule.claimedAmount += claimableAmount;
                totalTokenLocked[token] -= claimableAmount;
                IERC20(token).safeTransfer(beneficiary, claimableAmount);
            }
            remainingAmount = schedule.totalAmount - schedule.claimedAmount;
            if (remainingAmount > 0) {
                totalTokenLocked[token] -= remainingAmount;
                totalVestedAmount[token][beneficiary] -= remainingAmount;
                IERC20(token).safeTransfer(msg.sender, remainingAmount);
            }
        }
        schedule.revoked = true;
        emit VestingScheduleRevoked(token, beneficiary, scheduleId, remainingAmount);
    }

    function emergencyWithdrawToken(
        address token,
        uint256 amount
    ) external override onlyRole(ADMIN_ROLE) {
        uint256 balance = IERC20(token).balanceOf(address(this));
        if (amount > balance) revert InsufficientBalance();
        IERC20(token).safeTransfer(msg.sender, amount);
        emit EmergencyWithdraw(token, msg.sender, amount);
    }

    function pause() external  onlyRole(ADMIN_ROLE) {
        _pause();
    }

    function unpause() external  onlyRole(ADMIN_ROLE) {
        _unpause();
    }

    function _authorizeUpgrade(address  newImplementation) internal  override onlyRole(ADMIN_ROLE) {}
}
