#!/bin/bash

# Solidity and Go binding file paths and names
solidity_files=(
    # "cpuHeavy/CPUHeavy.sol"
    #"internalTxTC/InternalTxKIP17Token.sol"
    #"internalTxTC/InternalTxMainContract.sol"
    # "largeMemo/largeMemo.sol"
    # "readApiCallContract/ReadApiCallContract.sol"
    # "UserStorage/UserStorage.sol"
)

package_names=(
    # "cpuHeavyTC"
    #internalTxTC
    #internalTxTC
    # "largeMemoTC"
    # "readApiCallContractTC"
    # "UserStorageTC"
)

for i in "${!solidity_files[@]}"; do
    # Get the Solidity file path and package name
    solidity_file="${solidity_files[$i]}"
    package_name="${package_names[$i]}"

    # Get the directory and contract name
    dir=$(dirname "${solidity_file}")
    contract_name=$(basename "${solidity_file}" .sol)

    # Create build directory if it doesn't exist
    build_dir="${dir}/build"
    mkdir -p "${build_dir}"

    # Compile Solidity file
    solc --abi --bin --optimize --overwrite -o "${build_dir}" "${solidity_file}"

    # Generate Go binding file
    abigen --bin="${build_dir}/${contract_name}.bin" \
           --abi="${build_dir}/${contract_name}.abi" \
           --pkg="${package_name}" \
           --out="${dir}/${contract_name}.go"

    echo "Go binding file generated: ${dir}/${contract_name}.go with package ${package_name}"
done

# Cleanup: Remove build directories
for dir in "${solidity_files[@]}"; do
    dir_path=$(dirname "${dir}")
    build_dir="${dir_path}/build"
    if [ -d "${build_dir}" ]; then
        rm -rf "${build_dir}"
        echo "Removed build directory: ${build_dir}"
    fi
done
