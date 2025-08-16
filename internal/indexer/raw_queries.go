package indexer

const getBlocksGQL = /*graphql*/ `
query getBlocks($height_eq: Int!, $height_gt: Int!, $height_lt: Int!) {
    getBlocks(
        where: {
            _or: [
                { height: { eq: $height_eq } }
                { height: { gt: $height_gt, lt: $height_lt } }
            ]
        }
    ) {
        hash
        height
        time
        total_txs
        num_txs
    }
}`

const getTransactionsGQL = /*graphql*/ `
query getTransactions(
    $height_eq: Int!
    $height_gt: Int!
    $height_lt: Int!
    $index_eq: Int!
    $index_gt: Int!
    $index_lt: Int!
) {
    getTransactions(
        where: {
            _and: [
                {
                    _or: [
                        { block_height: { eq: $height_eq } }
                        { block_height: { gt: $height_gt, lt: $height_lt } }
                    ]
                }
                {
                    _or: [
                        { index: { eq: $index_eq } }
                        { index: { gt: $index_gt, lt: $index_lt } }
                    ]
                }
            ]
        }
    ) {
        index
        hash
        success
        block_height
        gas_wanted
        gas_used
        memo
        gas_fee {
            amount
            denom
        }
        messages {
            route
            typeUrl
            value {
                ... on BankMsgSend {
                    from_address
                    to_address
                    amount
                }
                ... on MsgAddPackage {
                    creator
                    send # deposit
                    package {
                        name
                        path
                        files {
                            name
                            body
                        }
                    }
                }
                ... on MsgCall {
                    pkg_path
                    func
                    send
                    caller
                    args
                }
                ... on MsgRun {
                    caller
                    send
                    package {
                        name
                        path
                        files {
                            name
                            body
                        }
                    }
                }
            }
        }
        response {
            log
            info
            error
            data
            events {
                ... on GnoEvent {
                    type
                    func
                    pkg_path
                    attrs {
                        key
                        value
                    }
                }
            }
        }
    }
}`
