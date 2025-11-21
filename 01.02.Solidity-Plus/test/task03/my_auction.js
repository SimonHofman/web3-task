const {ethers, deployments} = require("hardhat");
const {expect} = require("chai");

describe("Test auction", async function () {
    it("Should be ok", async function () {
        await main();
    });
})

async function main() {
    const [signer, buyer] = await ethers.getSigners();
    await deployments.fixture(["deployMyNftAuction"]);

    const myNftAuctionProxy = await deployments.get("MyNftAuctionProxy");
    const myNftAuction = await ethers.getContractAt(
        "MyNftAuction",
        myNftAuctionProxy.address
    );

    const MyERC20 = await ethers.getContractFactory("MyERC20");
    const myERC20 = await MyERC20.deploy();
    await myERC20.waitForDeployment();
    const UsdcAddress = await myERC20.getAddress();

    let tx = await myERC20.connect(signer).transfer(buyer, ethers.parseEther("10000"));
    await tx.wait();
    console.log("balance of buyer::", await myERC20.balanceOf(signer));
    console.log("balance of buyer::", await myERC20.balanceOf(buyer));

    const priceProvider = await ethers.getContractFactory("PriceProvider");

    const priceFeedEthDeploy = await priceProvider.deploy(ethers.parseEther("10000"));
    const priceFeedEth = await priceFeedEthDeploy.waitForDeployment();
    const priceFeedEthAddress = await priceFeedEth.getAddress();
    console.log("ethFeed:", priceFeedEthAddress);

    const priceFeedUSDCDeploy = await priceProvider.deploy(ethers.parseEther("1"));
    const priceFeedUSDC = await priceFeedUSDCDeploy.waitForDeployment();
    const priceFeedUSDCAddress = await priceFeedUSDC.getAddress();
    console.log("usdcFeed:", await priceFeedUSDCAddress);

    const token2Usd = [{
        token: ethers.ZeroAddress,
        priceFeed: priceFeedEthAddress
    }, {
        token: UsdcAddress,
        priceFeed: priceFeedUSDCAddress
    }];

    for (let i = 0; i < token2Usd.length; i++) {
        const {token, priceFeed} = token2Usd[i];
        await myNftAuction.setPriceFeed(token, priceFeed);
    }

    // 1. 部署ERC721 合约
    const MyERC721 = await ethers.getContractFactory("MyERC721");
    const myERC721 = await MyERC721.deploy();
    await myERC721.waitForDeployment();
    const myERC721Address = await myERC721.getAddress();
    console.log("myERC721Address::", myERC721Address);

    for (let i = 0; i < 10; i++) {
        await myERC721.mint(signer.address, i + 1);
    }

    const tokenId = 1;

    // 2. 给代理合约授权
    await myERC721.connect(signer).setApprovalForAll(myNftAuctionProxy.address, true);

    await myNftAuction.createAuction(
        10,
        ethers.parseEther("0.00000001"),
        myERC721Address,
        tokenId
    );
    const auction = await myNftAuction.auctions(0);

    console.log("创建拍卖成功::", auction);

    // 3. 购买者参与拍卖
    tx = await myNftAuction.connect(buyer).placeBid(0, 0, ethers.ZeroAddress, {value: ethers.parseEther("0.00000002")});
    await tx.wait();

    // USDC参与竞价
    tx = await myERC20.connect(buyer).approve(myNftAuctionProxy.address, ethers.MaxUint256);
    await tx.wait();
    tx = await myNftAuction.connect(buyer).placeBid(0, ethers.parseEther("101"), UsdcAddress);
    result = await tx.wait();

    // 4. 结束拍卖
    await new Promise((resolve) => setTimeout(resolve, 10 * 1000));
    await myNftAuction.connect(signer).endAuction(0);

    // 验证结果
    const auctionResult = await myNftAuction.auctions(0);
    console.log("结束拍卖后读取拍卖成功::", auctionResult);
    expect(auctionResult.highestBidder).to.equal(buyer.address);
    expect(auctionResult.highestBid).to.equal(ethers.parseEther("101"));

    // 验证NFT所有权
    const owner = await myERC721.ownerOf(tokenId);
    console.log("owner::", owner);
    expect(owner).to.equal(buyer.address);
}

main()