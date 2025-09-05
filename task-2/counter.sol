// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Counter {
    uint256 private _count;
    address public owner;

    event CountIncremented(uint256 newCount);
    event CountReset();

    constructor() {
        owner = msg.sender;
        _count = 0;
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    function getCount() public view returns (uint256) {
        return _count;
    }

    function increment() public {
        _count += 1;
        emit CountIncremented(_count);
    }

    function reset() public onlyOwner {
        _count = 0;
        emit CountReset();
    }

    function setCount(uint256 newCount) public onlyOwner {
        _count = newCount;
        emit CountIncremented(_count);
    }
}