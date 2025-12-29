const {ethers, upgrades} = require("hardhat");
const esDex_name = "EasySwapOrderBook";
const esDex_address = "0xcEE5AA84032D4a53a0F9d2c33F36701c3eAD5895"

const esVault_name = "EasySwapVault";
const esVault_address = "0xaD65f3dEac0Fa9Af4eeDC96E95574AEaba6A2834"
async function main() {
    const [signer, owner] = await ethers.getSigners();
    console.log(signer.address, ": signer");

    let exDex = await ethers.getContractFactory("EasySwapOrderBook");
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });