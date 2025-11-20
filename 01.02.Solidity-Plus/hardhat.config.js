require("@nomicfoundation/hardhat-toolbox");
require("@openzeppelin/hardhat-upgrades");
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


task("accounts", "Prints the list of accounts", async (taskArgs, hre) => {
    const accounts = await hre.ethers.getSigners();
    accounts.forEach((account, index) => {
        console.log(`账户 ${index}: ${account.address}`);
    });
});
