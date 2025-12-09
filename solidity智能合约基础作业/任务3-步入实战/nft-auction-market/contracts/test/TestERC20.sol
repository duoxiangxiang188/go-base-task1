//SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract TestERC20 is ERC20,Ownable {
    constructor() ERC20("TestToken","TST") Ownable(msg.sender) {
        _mint(msg.sender,1000000 *10**decimals());
    }
    function mint(address to ,uint256 amount) public onlyOwner{
        _mint(to,amount);
    }
    // 为了兼容测试，暴露decimals（ERC20默认18位）
    function decimals() public pure override returns (uint8) {
        return 18;
    }

}