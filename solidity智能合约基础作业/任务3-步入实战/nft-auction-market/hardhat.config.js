// hardhat.config.js
// 仅保留一个verify插件：@nomicfoundation/hardhat-verify
require("dotenv").config();
require("@nomicfoundation/hardhat-chai-matchers");
require("@nomicfoundation/hardhat-ethers"); 
require("@nomicfoundation/hardhat-verify");
require("hardhat-deploy");
require("hardhat-gas-reporter");
require("solidity-coverage");
require("@openzeppelin/hardhat-upgrades");



const PRIVATE_KEY = process.env.PRIVATE_KEY ;
const SEPOLIA_RPC_URL = process.env.SEPOLIA_RPC_URL ;
const ETHERSCAN_API_KEY = process.env.ETHERSCAN_API_KEY || ""; // 验证合约用

module.exports = {
  solidity: {
    compilers: [{ version: "0.8.28" }],
  },
  networks: {
    hardhat: {},
    sepolia: {
      url: SEPOLIA_RPC_URL,
      accounts: [PRIVATE_KEY],
      saveDeployments: true,
      chainId: 11155111,
    },
  },
  namedAccounts: {
    deployer: { default: 0 },
  },
  gasReporter: {
    enabled: true,
    currency: "USD",
    outputFile: "gas-report.txt",
  },
  // 适配@nomicfoundation/hardhat-verify的验证配置
  verify: {
    etherscan: {
      apiKey: ETHERSCAN_API_KEY,
    },
  },
  mocha: {
    timeout: 200000,
  },
  // 强制适配ethers v5
  ethers: {
    version: "5.7.2",
  },
};