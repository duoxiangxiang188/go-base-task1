// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract BeggingContract {
    address private immutable _owner;

    mapping (address => uint256) private _donations;

    address[] private _donors;

    // uint256 public donationStartTime;
    // uint256 public donationEndTime;


    event DonationReceived(address indexed donor, uint256 amount, uint256 timestamp);
    event Withdrawal(address indexed owner, uint256 amount, uint256 timestamp);

    modifier onlyOwner(){
        require(msg.sender == _owner,"Only owner can call this function");
        _;
    }
    // modifier onlyDuringDonationPeriod(){
    //     require(block.timestamp >= donationStartTime &&block.timestamp <=donationEndTime,"Donation is not allowed at this time");
    //     _;
    // }
    // constructor(uint256 start,uint256 end){
    //     _owner = msg.sender;
    //     donationStartTime=start;
    //     donationEndTime=end;

    // }
    constructor(){
        _owner = msg.sender;

    }
    function donate()external  payable {
        require(msg.value>0,"Donation amount must be greater than 0");

        if(_donations[msg.sender] == 0){
            _donors.push(msg.sender);

        }
        
        _donations[msg.sender] +=msg.value;

        emit DonationReceived(msg.sender, msg.value, block.timestamp);

    }

    function withdraw()external  onlyOwner{
        uint256 balance = address(this).balance ;
        
        require(balance>0,"No funds to withdraw");

        (bool success, )=_owner.call{value:balance}("");
        require(success,"Withdrawal failed");
        emit  Withdrawal(_owner, balance, block.timestamp);


    }
    function getDonation(address donor) external view returns (uint256){
        return _donations[donor];
    }

    function getTop3Donors()external view returns(address[] memory,uint256[] memory) {
        address[] memory topAddresses = new address[](3);
        uint256[] memory topAmounts = new uint256[](3);

        address[] memory temDonors = _donors;
        uint256 length = temDonors.length;
        for(uint256 i = 0;i<length;i++){
            for (uint256 j=i+1;j<length;j++){
                if(_donations[temDonors[j]]>_donations[temDonors[i]]){
                    (temDonors[i],temDonors[j])=(temDonors[j],temDonors[i]);                    
                }
            }
        }
        for(uint256 i=0;i<3;i++){
            if(i<length){
                topAddresses[i]=temDonors[i];
                topAmounts[i]=_donations[temDonors[i]];
            }else {
                topAddresses[i]=address(0);
                topAmounts[i]=0;
            }
        }
        return (topAddresses,topAmounts);
    }

    function getContractBalance()external view returns(uint256) {
        return address(this).balance;

    }

    function getOwner()external view returns(address) {
        return _owner;
    }


}