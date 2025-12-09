// contracts/mocks/MockPriceFeed.sol
pragma solidity ^0.8.20;

import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

contract MockPriceFeed is AggregatorV3Interface {
    int256 public immutable price; // 固定价格（如 1e8 = 1 USD/代币）
    uint8 public immutable _mockDecimals; // 更独特的变量名，避免冲突
    string public constant _mockDescription = "Mock TestERC20/USD Price Feed";
    uint256 public constant _mockVersion = 1;

    constructor(int256 _price, uint8 decimals) {
        price = _price;
        _mockDecimals = decimals; // 同步修改赋值
    }

    // 实现接口方法：返回精度
    function decimals() external view returns (uint8) {
        return _mockDecimals; // 同步修改返回值
    }

    // 其余方法（description()、version()等）同理，同步修改变量名
    function description() external view returns (string memory) {
        return _mockDescription;
    }

    function version() external view returns (uint256) {
        return _mockVersion;
    }

    // 其余方法（getRoundData、latestRoundData）保持不变
    // 实现接口方法：返回指定轮次的价格
    function getRoundData(uint80 _roundId)
        external
        view
        returns (
            uint80 roundId,
            int256 answer,
            uint256 startedAt,
            uint256 updatedAt,
            uint80 answeredInRound
        )
    {
        return (
            _roundId,
            price,
            block.timestamp,
            block.timestamp,
            _roundId
        );
    }

    // 实现接口方法：返回最新轮次价格
    function latestRoundData()
        external
        view
        returns (
            uint80 roundId,
            int256 answer,
            uint256 startedAt,
            uint256 updatedAt,
            uint80 answeredInRound
        )
    {
        return (
            1,
            price,
            block.timestamp,
            block.timestamp,
            1
        );
    }
}