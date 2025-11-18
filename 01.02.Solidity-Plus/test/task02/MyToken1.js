const {config} = require("chai");
const hre = require("hardhat");

async function getBalance(contract, address) {
    const balance = await contract.balanceOf(address);
    console.log(`${address} balance is :`, hre.ethers.formatEther(balance));
}

async function transferERC20(contract, address, amount) {
    const tx = await contract.transfer(address, amount);
    await tx.wait();
    console.log(`Transfer transaction hash:`, tx.hash);
}

async function approveERC20(contract, address, amount) {
    const tx = await contract.approve(address, amount);
    await tx.wait();
    console.log(`Approve transaction hash:`, tx.hash);
}

async function getAllowance(contract, owner, spender) {
    const allowance = await contract.allowance(owner, spender);
    console.log(`${owner} has allowance of ${spender}:`, allowance);
}

async function transferFrom(contract, from, spender, to, amount){
    // 获取用户的签名者
    const userSigner = await hre.ethers.getSigner(spender);
    // 连接到合约
    const contractAsUser = contract.connect(userSigner);
    const tx = await contractAsUser.transferFrom(from, to, amount);
    await tx.wait();
    console.log(`transferFrom transaction hash:`, tx.hash);
}

async function main() {
    const {deployer, user1, user2, user3, user4} = await hre.getNamedAccounts()
    console.log("deployer account:", deployer);
    console.log("user1 account:", user1);
    console.log("user2 account:", user2);
    console.log("user3 account:", user3);
    console.log("user4 account:", user4);

    const MyToken1 = await hre.ethers.getContractFactory("MyERC20Token1");
    const contract = await MyToken1.deploy("MT", "MT");
    await contract.waitForDeployment();
    console.log("Contract is deployed to:", await contract.getAddress());

    result = await contract.name();
    console.log("name is : ", result);
    result = await contract.symbol();
    console.log("symbol is : ", result);
    result = await contract.decimals();
    console.log("decimals is : ", result);
    result = await contract.totalSupply();
    console.log("totalSupply is : ", result);

    await getBalance(contract, deployer);
    await getBalance(contract, user1);
    await getBalance(contract, user2);
    await getBalance(contract, user3);
    await getBalance(contract, user4);

    await transferERC20(contract, user1, BigInt(10 * 10 ** 18));
    await transferERC20(contract, user2, BigInt(20 * 10 ** 18));
    await transferERC20(contract, user3, BigInt(30 * 10 ** 18));
    await transferERC20(contract, user4, BigInt(40 * 10 ** 18));

    await getBalance(contract, deployer);
    await getBalance(contract, user1);
    await getBalance(contract, user2);
    await getBalance(contract, user3);
    await getBalance(contract, user4);

    await approveERC20(contract, user1, BigInt(100 * 10 ** 18));
    await getAllowance(contract, deployer, user1);

    await transferFrom(contract, deployer, user1, user4, BigInt(10 * 10 ** 18))
    await getBalance(contract, deployer);
    await getBalance(contract, user1);
    await getBalance(contract, user2);
    await getBalance(contract, user3);
    await getBalance(contract, user4);
}

main()