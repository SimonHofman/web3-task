const {config} = require("chai");
const hre = require("hardhat");

async function donateFromUser(contract, userAddress, amountEth) {
    // 获取用户的签名者
    const userSigner = await hre.ethers.getSigner(userAddress);
    // 连接到合约
    const contractAsUser = contract.connect(userSigner);
    // 将ETH金额转换为wei
    const donationAmount = hre.ethers.parseEther(amountEth.toString());
    // 调用 donate 函数并发送ETH
    const tx = await contractAsUser.donate({
        value: donationAmount
    });
    await tx.wait(); // 等待交易确认
    console.log(`${userAddress} donated ${amountEth} ETH to the contract.`)
    console.log("Donation transaction hash:", tx.hash);
}

async function getDonation(contract, address) {
    const donationAmount = await contract.getDonation(address);
    console.log(`${address} has donated:`, hre.ethers.formatEther(donationAmount), "ETH");
}

async function getBalance(contract, address) {
    const balance = await hre.ethers.provider.getBalance(address);
    console.log(`${address} balance is :`, hre.ethers.formatEther(balance)) ;
}

async function changeOwner(contract, newOwner) {
    const tx = await contract.transferOwnership(newOwner);
    await tx.wait(); // 等待交易确认
    console.log("Change owner transaction hash:", tx.hash);
}

async function getContractOwner(contract) {
    const result = await contract.getOwner();
    console.log("Owner:", result);
    return result;
}

async function ownerWithDraw(contract, user) {
    // 获取用户的签名者
    const userSigner = await hre.ethers.getSigner(user);
    // 连接到合约
    const contractAsUser = contract.connect(userSigner);
    await contractAsUser.withdraw();
}

async function getTop3(contract) {
    const top3 = await contract.getTop3Donors();
    console.log("Top 3:", top3);
}

async function main() {
    const {deployer, user1, user2, user3, user4} = await hre.getNamedAccounts();
    console.log("deployer account:", deployer);
    console.log("user1 account:", user1);
    console.log("user2 account:", user2);
    console.log("user3 account:", user3);
    console.log("user4 account:", user4);

    const BeggingContractExtra = await hre.ethers.getContractFactory("BeggingContractExtra");
    const contract = await BeggingContractExtra.deploy(1763078700, 1763478700);
    await contract.waitForDeployment();
    console.log("Contract is deployed to:", await contract.getAddress())

    await donateFromUser(contract, user1, 1)
    await getBalance(user1)
    await donateFromUser(contract, user2, 1.5)
    await getBalance(user2)
    await donateFromUser(contract, user3, 2)
    await getBalance(user3)

    await getTop3(contract);

    await getDonation(contract, user1);
    await getDonation(contract, user2);
    await getDonation(contract, user3);
    await getDonation(contract, user4);

    await donateFromUser(contract, user1, 1)
    await getTop3(contract);

    // await getContractOwner(contract);
    await changeOwner(contract, user4)

    await ownerWithDraw(contract, user4);
    await getBalance(user4)
}

main()