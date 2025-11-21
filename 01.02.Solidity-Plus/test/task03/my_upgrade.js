const {ethers, deployments, upgrades} = require("hardhat");
const {expect} = require("chai");

describe("Test upgrade", async function () {
    it("Should be able to deploy", async function () {
        const [signer, buyer] = await ethers.getSigners();

        // 1. 部署业务合约
        await deployments.fixture(["deployMyNftAuction"]);
        const myNftAuctionProxy = await deployments.get("MyNftAuctionProxy");
        console.log(myNftAuctionProxy);

        // 2. 部署 ERC721 合约
        const MyERC721 = await ethers.getContractFactory("MyERC721");
        const myERC721 = await MyERC721.deploy();
        await myERC721.waitForDeployment();
        const myERC721Address = await myERC721.getAddress();
        console.log("myERC721Address::", myERC721Address);

        // mint 10个 NFT
        for (let i = 0; i < 10; i++) {
            await myERC721.mint(signer.address, i + 1);
        }

        const tokenId = 1;

        // 给代理合约授权
        await myERC721.connect(signer).setApprovalForAll(myNftAuctionProxy.address, true);

        // 3. 调用 createAuction 方法创建拍卖
        const myNftAuction = await ethers.getContractAt(
            "MyNftAuction",
            myNftAuctionProxy.address
        );

        await myNftAuction.createAuction(
            100 * 1000,
            ethers.parseEther("0.01"),
            myERC721Address,
            1
        );

        const auction = await myNftAuction.auctions(0);
        console.log("创建拍卖成功::", auction);

        const implAddress1 = await upgrades.erc1967.getImplementationAddress(
            myNftAuctionProxy.address
        );

        // 4. 升级合约
        await deployments.fixture(["upgradeMyNftAuction"]);
        const implAddress2 = await upgrades.erc1967.getImplementationAddress(
            myNftAuctionProxy.address
        );

        // 4. 读取合约的 auction[0]
        const auction2 = await myNftAuction.auctions(0);
        console.log("升级后读取拍卖成功::", auction2);

        console.log("implAddress1::", implAddress1, "\nimplAddress2::", implAddress2);

        const myNftAuctionV2 = await ethers.getContractAt(
            "MyNftAuctionV2",
            myNftAuctionProxy.address
        );
        const hello = await myNftAuctionV2.testHello();
        console.log("hello:", hello);

        expect(auction2.startTime).to.equal(auction.startTime);
    });
});