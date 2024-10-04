// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract Verifier {
    struct Node {
        address id;
        string ip;
    }

    mapping(address => Node) nodes;

    function add(string calldata _ip) public {
        nodes[msg.sender] = Node({id: msg.sender, ip: _ip});
    }

    function isAdded(string calldata ip) public view returns (bool) {
        return keccak256(bytes(nodes[msg.sender].ip)) == keccak256(bytes(ip));
    }

    function verify(
        address _addr,
        string calldata _ip
    ) public view returns (bool) {
        if (keccak256(bytes(nodes[_addr].ip)) == keccak256(bytes(_ip))) {
            return true;
        }

        return false;
    }
}
