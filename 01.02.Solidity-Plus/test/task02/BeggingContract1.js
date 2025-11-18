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

async function getBalance(address) {
    const balance = await hre.ethers.provider.getBalance(address);
    console.log(`${address} balance is :`, hre.ethers.formatEther(balance)) ;
}
async function main() {
    const {deployer, user1, user2, user3, user4} = await hre.getNamedAccounts();
    console.log("deployer account:", deployer);
    console.log("user1 account:", user1);
    console.log("user2 account:", user2);
    console.log("user3 account:", user3);
    console.log("user4 account:", user4);

    const BeggingContract1 = await hre.ethers.getContractFactory("BeggingContract");
    const contract = await BeggingContract1.deploy();
    await contract.waitForDeployment();
    console.log("Contract is deployed to:", await contract.getAddress())

    await donateFromUser(contract, user1, 1)
    await getBalance(user1)
    await donateFromUser(contract, user2, 1.5)
    await getBalance(user2)
    await donateFromUser(contract, user3, 2)
    await getBalance(user3)

    await getDonation(contract, user1);
    await getDonation(contract, user2);
    await getDonation(contract, user3);
    await getDonation(contract, user4);

    await contract.withdraw();
    await getBalance(deployer)
}

main()