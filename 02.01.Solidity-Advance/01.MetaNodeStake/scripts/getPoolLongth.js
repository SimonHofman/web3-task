const {ethers} = require("hardhat");

async function main() {
    const MetaNodeStake = await ethers.getContractAt("MetaNodeStake", "0x0165878A594ca255338adfa4d48449f69242Eb8F");

    console.log(await MetaNodeStake.poolLength());
}

main()