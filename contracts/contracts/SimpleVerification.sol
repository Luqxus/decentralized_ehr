// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contract/ownership/Ownable.sol";
import "@openzeppelin/contracts/utils/cryptography/MerkleProof.sol";

contract SimpleVerifier Ownable {
	bytes32 private root;

	constructor(bytes32 _root) {
		root = _root;
	}

	function verify(
		[]bytes32 calldata proof,
		string calldata hash,
	) public pure returns(bool) {
		bytes32 leaf = keccak256(bytes(hash));
		return MerkleProof.verify(proof, root, leaf)
	}
}
