// SPDX-License-Identifier: MIT
pragma solidity ^0.8;
import "@openzeppelin/contracts/access/Ownable.sol"; // 引入Ownable控制权限
contract Voting is Ownable {
   
    mapping(bytes32 => uint256) private _candidateVoteCount;
    mapping(bytes32 => string) private _candidateNames;
    struct CandidateVoteWithName  {
        string candidateName;
        uint256 votes;
    }

    bytes32[] private _candidates;

    constructor() Ownable(msg.sender) {}
    function vote(bytes32 candidate)  public {
        _candidateVoteCount[candidate] += 1;
        if (!_isCandidateExist(candidate)) {
            _candidates.push(candidate);
            
        }
    }
    function getVotes(bytes32 candidate) public  view returns (uint) {
        return _candidateVoteCount[candidate];
    }

    function resetVotes() external  onlyOwner {
        for (uint256 i = 0; i < _candidates.length; i++) {
            _candidateVoteCount[_candidates[i]] = 0;
          
        }
        delete _candidates;
        
    }

    function _isCandidateExist(bytes32 candidate) internal view returns (bool) {
        for (uint256 i = 0; i < _candidates.length; i++) {
            if (_candidates[i] == candidate) {
                return true;
            }
        }
        return false;
    }

    function voteWithString(string calldata candidate) external {
        bytes32 candidateHash = keccak256(bytes(candidate));
        vote(candidateHash);

        if(bytes(_candidateNames[candidateHash]).length == 0) {
            _candidateNames[candidateHash] = candidate;
        }
    }
      function getAllVotes() external view returns (CandidateVoteWithName[] memory ) {

            CandidateVoteWithName[] memory result = new CandidateVoteWithName[](_candidates.length);

            for(uint256 i = 0; i < _candidates.length; i++ ) {
                bytes32 candidate = _candidates[i];
                result[i] = CandidateVoteWithName({
                    candidateName: _candidateNames[candidate],
                    votes: _candidateVoteCount[candidate]
                });
            }
        return result ;
    }

    function getVotesWithString(string calldata candidate) external view returns (uint256) {
        return getVotes(keccak256(abi.encodePacked(candidate)));
    }
}