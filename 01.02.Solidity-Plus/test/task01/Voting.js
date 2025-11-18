const {config} = require("chai");
const hre = require("hardhat");

async function main() {
    const {deployer, user1, user2, user3, user4} = await hre.getNamedAccounts();
    console.log("deployer account:", deployer);
    console.log("user1 account:", user1);
    console.log("user2 account:", user2);
    console.log("user3 account:", user3);
    console.log("user4 account:", user4);

    const Voting = await hre.ethers.getContractFactory("Voting");
    const contract = await Voting.deploy();
    await contract.waitForDeployment();
    console.log("Contract is deployed to:", await contract.getAddress())

    await contract.vote(user1);
    await contract.vote(user1);
    await contract.vote(user1);
    await contract.vote(user1);
    await contract.vote(user1);

    await contract.vote(user2);
    await contract.vote(user2);
    await contract.vote(user2);
    await contract.vote(user2);

    try {
        const result = await contract.getVotes(user1);
        console.log("user1 has ", result)
    } catch (error) {
        console.error("Error:", error);
    }

    await contract.resetVotes();

    try {
        const result = await contract.getVotes(user1);
        console.log("user1 has ", result)
    } catch (error) {
        console.error("Error:", error);
    }
}
main()