// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

interface IERC20 {
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner,address indexed spender, uint256 value);

    function balanceOf(address account) external  view returns (uint256);
    function transfer(address to, uint256 value) external returns (bool);
    function approve(address spender, uint256 value) external returns (bool);
    function transferFrom(address from, address to, uint256 value)external returns (bool);
    function totalSupply() external view returns (uint256);
    function allowance(address owner, address spender) external view returns (uint256);
    
}

contract SimpleERC20 is IERC20 {
    string public  name = "SimpleERC20";
    string public symbol= "SIM";
    uint8 public decimals = 18;
    mapping (address => uint256)private  _balances;
    mapping (address => mapping (address=>uint256)) private _allowances;
    uint256 private _totalSupply;
    address private _owner;

    constructor() {
        _owner = msg.sender;
    }
    modifier onlyOwner(){
        require(msg.sender == _owner, "Only owner can call");
        _;
    }


    function balanceOf(address account) external   view  override returns (uint256){
        return _balances[account];
    }

    function transfer(address to, uint256 value) external   override returns (bool){
        address from = msg.sender;
        require(_balances[from] >= value, "Insufficient balance");
        require(to != address(0), "Transfer to zero address");

        _balances[from] -= value;
        _balances[to] += value;
        emit Transfer(from, to, value);
        return true;
    }

    function approve(address spender, uint256 value) external  override returns (bool){
        address owner =msg.sender;
        require(spender != address(0), "Approve to zero address");

        _allowances[owner][spender] = value;
        emit Approval(owner, spender, value);
        return true;

    }

    function transferFrom(address from, address to, uint256 value) external override returns (bool){
        address spender = msg.sender;
        uint256 allowed = _allowances[from][spender];
        require(_balances[from] >= value, "Insufficient balance");
        require(allowed >= value, "Allowance exceeded");
        require(to != address(0), "Transfer to zero address");

        _allowances[from][spender] -= value;
        _balances[from] -= value;
        _balances[to] += value;
        emit Transfer(from, to, value);
        return true;

    }

    function mint(address to, uint256 value) public onlyOwner returns (bool){
        require(to!= address(0), "Mint to zero address");

        _totalSupply += value;
        _balances[to] += value;

        emit Transfer(address(0), to, value);
        return true;

    }

    function totalSupply()external  view override  returns (uint256){
        return _totalSupply;

    }
    function allowance(address owner, address spender)external  view override  returns (uint256){
        return _allowances[owner][spender];
    }
}