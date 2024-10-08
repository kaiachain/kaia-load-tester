// Derived from BlockBench's SmallBank benchmark.
pragma solidity ^0.8.24;

contract SmallBank {
    mapping(string=>uint) savingStore;
    mapping(string=>uint) checkingStore;

    function almagate(string memory arg0, string memory arg1) public {
       uint bal1 = savingStore[arg0];
       uint bal2 = checkingStore[arg1];

       checkingStore[arg0] = 0;
       savingStore[arg1] = bal1 + bal2;
    }

    function getBalance(string memory arg0) public view returns (uint balance) {
        uint bal1 = savingStore[arg0];
        uint bal2 = checkingStore[arg0];

        balance = bal1 + bal2;
        return balance;
    }

    function updateBalance(string memory arg0, uint arg1) public {
        uint bal1 = checkingStore[arg0];
        uint bal2 = arg1;

        checkingStore[arg0] = bal1 + bal2;
    }

    function updateSaving(string memory arg0, uint arg1) public {
        uint bal1 = savingStore[arg0];
        uint bal2 = arg1;

        savingStore[arg0] = bal1 + bal2;
    }

    function sendPayment(string memory arg0, string memory arg1, uint arg2) public {
        uint bal1 = checkingStore[arg0];
        uint bal2 = checkingStore[arg1];
        uint amount = arg2;

        bal1 -= amount;
        bal2 += amount;

        checkingStore[arg0] = bal1;
        checkingStore[arg1] = bal2;
    }

    function writeCheck(string memory arg0, uint arg1) public {
        uint bal1 = checkingStore[arg0];
        uint bal2 = savingStore[arg0];
        uint amount = arg1;

        if (amount < bal1 + bal2) {
            checkingStore[arg0] = bal1 - amount - 1;
        }
        else {
            checkingStore[arg0] = bal1 - amount;
        }
    }
}
