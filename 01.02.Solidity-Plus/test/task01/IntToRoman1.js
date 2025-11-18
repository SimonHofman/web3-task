const {config} = require("chai");
const {ethers} = require("hardhat")

async function main() {
    const [deployer] = await ethers.getSigners();
    console.log("Deploying with:", deployer.address);

    const IntToRoman1 = await ethers.getContractFactory("IntToRoman1");
    const contract = await IntToRoman1.deploy();
    await contract.waitForDeployment();
    console.log("Contract is deployed to:", await contract.getAddress());

    const intArray = [3749, 58, 1994]
    const strArray = ["MMMDCCXLIX", "LVIII", "MCMXCIV"]

    for (let i = 0; i < intArray.length; i++) {
        try {
            const result = await contract.intToRoman(intArray[i]);
            console.log("Result:", result);
            console.log("Expected:", strArray[i]);
            console.log("Is correct:", result === strArray[i]);
        } catch (error) {
            console.error("Error:", error);
        }
    }
}
main()