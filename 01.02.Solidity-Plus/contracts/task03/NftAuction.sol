// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

import "hardhat/console.sol";

contract NftAuction is Initializable, UUPSUpgradeable {
    // 结构体
    struct Auction {
        // 卖家
        address seller;
        // 开始时间
        uint256 startTime;
        // 拍卖持续时间
        uint256 duration;
        // 初始价格
        uint256 startPrice;
        // 是否开始
        bool ended;
        // 最高出价者
        address highestBidder;
        // 最高出价
        uint256 highestBid;
        // NFT合约地址
        address nftContract;
        // NFT ID
        uint256 tokenId;
        // 参与竞价的资产类型：0x地址代表eth，其他地址代表erc20
        address tokenAddress;
    }

    // 拍卖的状态
    mapping(uint256 => Auction) public auctions;
    // 下一个拍卖ID
    uint256 public nextAuctionId;
    // 管理员地址
    address public admin;

    // 预言机
    mapping(address => AggregatorV3Interface) public priceFeeds;

    function initialize() public initializer {
        admin = msg.sender;
    }

    function setPriceFeed(address tokenAddress, address _priceFeed) public {
        priceFeeds[tokenAddress] =AggregatorV3Interface(_priceFeed);
    }

    function getChainlinkDataFeedLatestAnswer(address tokenAddress) public view returns (int) {
        AggregatorV3Interface priceFeed = priceFeeds[tokenAddress];
        (
        /*uint80 roundID*/,
            int256 price,
        /*uint startedAt*/,
        /*uint timeStamp*/,
        /*uint80 answeredInRound*/
        ) = priceFeed.latestRoundData();
        return price;
    }

    // 创建拍卖
    function createAuction(
        uint256 _duration,
        uint256 _startPrice,
        address _nftContractAddress,
        uint256 _tokenId
    ) public {
        require(msg.sender == admin, "Only admin can create auctions");
        // 检查参数
        require(_duration > 0, "Duration must be greater than 0");
        require(_startPrice > 0, "Start price must be greater than 0");

        // 转移NFT到合约
        IERC721(_nftContractAddress).safeTransferFrom(msg.sender, address(this), _tokenId);

        // 创建新的拍卖记录并存储到autions映射中
        auctions[nextAuctionId] = Auction({
            seller: msg.sender,
            startTime: block.timestamp,
            duration: _duration,
            startPrice: _startPrice,
            ended: false,
            highestBidder: address(0),
            highestBid: 0,
            nftContract: _nftContractAddress,
            tokenId: _tokenId,
            tokenAddress: address(0)
        });

        // 拍卖ID加1
        nextAuctionId++;
    }

    // 买家参与竞拍
    function placeBid(
        uint256 _auctionId,
        uint256 amount,
        address _tokenAddress
    ) external  payable{
        // 判断拍卖状态
        Auction storage auction = auctions[_auctionId];
        require(
            ! auction.ended && block.timestamp < auction.startTime + auction.duration,
            "Auction has ended"
        );

        // 判断出价是否大于当前最高价
        uint _payment;
        if (_tokenAddress != address(0)) {
            _payment = amount * uint(getChainlinkDataFeedLatestAnswer(_tokenAddress));
        } else {
            amount = msg.value;
            _payment = amount * uint(getChainlinkDataFeedLatestAnswer(address(0)));
        }

        uint startPriceValue = auction.startPrice *
                            uint(getChainlinkDataFeedLatestAnswer(auction.tokenAddress));
        uint highestBidValue = auction.highestBid *
                            uint(getChainlinkDataFeedLatestAnswer(auction.tokenAddress));

        require(
            _payment > startPriceValue && _payment > highestBidValue,
            "Bid must higher the the current highest bid"
        );

        // 退还前最高价
        if (auction.highestBid >0 ) {
            if (auction.tokenAddress == address(0)) {
                payable(auction.highestBidder).transfer(auction.highestBid);
            } else {
                IERC20(auction.tokenAddress).transfer(auction.highestBidder, auction.highestBid);
            }
        }

        auction.tokenAddress = _tokenAddress;
        auction.highestBidder = msg.sender;
        auction.highestBid = amount;
    }

    // 结束拍卖
    function endAuction(uint256 _auctionId) external payable {
        Auction storage auction = auctions[_auctionId];
        console.log(
            "endAuction",
            auction.startTime,
            auction.duration,
            block.timestamp
        );

        // 判断当前拍卖是否结束
        require(
            ! auction.ended && block.timestamp >= auction.startTime + auction.duration,
            "Auction has not ended"
        );

        // 转移NFT到最高出价者
        IERC721(auction.nftContract).safeTransferFrom(
            address(this),
            auction.highestBidder,
            auction.tokenId
        );

        // 转移剩余资金到卖家
        if (auction.highestBid > 0) {
            if (auction.tokenAddress == address(0)) {
                payable(auction.seller).transfer(auction.highestBid);
            } else {
                IERC20(auction.tokenAddress).transfer(auction.seller, auction.highestBid);
            }
        }
        auction.ended = true;
    }

    function _authorizeUpgrade(address) internal view override {
        require(msg.sender == admin, "Only admin can upgrade");
    }

    function onERC721Received(
        address operator,
        address from,
        uint256 tokenId,
        bytes calldata data
    ) external pure returns (bytes4){
        return this.onERC721Received.selector;
    }
}
