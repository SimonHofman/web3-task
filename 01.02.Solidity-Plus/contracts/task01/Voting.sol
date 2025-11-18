// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

contract Voting {
    address[] private candidates;
    mapping(address => uint) votes;

    function vote(address _candidate) public {
        votes[_candidate] += 1;
        candidates.push(_candidate);
    }

    function getVotes(address _candidate) public view returns (uint) {
        return votes[_candidate];
    }

    function resetVotes() public {
        for (uint i = 0; i < candidates.length; i++) {
            votes[candidates[i]] = 0;
        }
        delete candidates;
    }
}
