require("@nomicfoundation/hardhat-toolbox");
require('hardhat-deploy');

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.28",
    namedAccounts: {
        deployer: 0,
        user1: 1,
        user2: 2,
        user3: 3,
        user4: 4,
    }
};
