pragma solidity ^0.5.0;

import "./IKIP17.sol";

/**
 * @title KIP-17 Non-Fungible Token Standard, optional enumeration extension
 * @dev See https://kips.kaia.io/KIPs/kip-17
 */
contract IKIP17Enumerable is IKIP17 {
    function totalSupply() public view returns (uint256);
    function tokenOfOwnerByIndex(address owner, uint256 index) public view returns (uint256 tokenId);

    function tokenByIndex(uint256 index) public view returns (uint256);
}
