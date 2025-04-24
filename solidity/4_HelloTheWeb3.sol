// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;




contract HelloTheWeb3 {
    string public tw="Hello World";
    uint256 public version=1;


    constructor(){
        version=2;
    }
    function getTw()external view returns (string memory){
       return tw;
    }
    function setTw(string calldata _tw) external {
         tw=_tw;
     }
    function getVersion()external  view returns (uint256){
       return version;
    }
    function setVersion(uint256  _version) external {
         version=_version;
     }
}