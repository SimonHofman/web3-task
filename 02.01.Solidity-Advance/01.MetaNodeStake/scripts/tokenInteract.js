const {ethers} = require("hardhat")

async function main() {
    const stakeContract = ethers.getContractAt("MetaNodeStake", "0x0165878A594ca255338adfa4d48449f69242Eb8F")
    const data = await stakeContract.MetaNode()
    console.log(data);
}

main()