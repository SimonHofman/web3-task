require("@nomicfoundation/hardhat-toolbox");
require("@openzeppelin/hardhat-upgrades");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
    solidity: {
        version: "0.8.28",
        settings: {
            optimizer: {
                enabled: true,
                runs: 20000,  // 增加优化运行次数
            },
        },
    },
    networks: {
        hardhat: {
            allowUnlimitedContractSize: true,  // 在开发网络中启用
        },
    },
};
