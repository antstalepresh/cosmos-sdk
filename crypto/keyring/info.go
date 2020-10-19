package keyring

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/types"
)

// Info is the publicly exposed information about a keypair
type Info interface {
	// Human-readable type for key listing
	GetType() KeyType
	// Name of the key
	GetName() string
	// Public key
	GetPubKey() crypto.PubKey
	// Address
	GetAddress() types.AccAddress
	// WithNameAndAddress
	WithNameAndAddress(string, types.AccAddress) Info
	// Bip44 Path
	GetPath() (*hd.BIP44Params, error)
	// Algo
	GetAlgo() hd.PubKeyType
}

var (
	_ Info = &localInfo{}
	_ Info = &ledgerInfo{}
	_ Info = &offlineInfo{}
	_ Info = &multiInfo{}
)

// localInfo is the public information about a locally stored key
// Note: Algo must be last field in struct for backwards amino compatibility
type localInfo struct {
	Name         string           `json:"name"`
	PubKey       crypto.PubKey    `json:"pubkey"`
	PrivKeyArmor string           `json:"privkey.armor"`
	Algo         hd.PubKeyType    `json:"algo"`
	Address      types.AccAddress `json:"address"`
}

func newLocalInfo(name string, pub crypto.PubKey, privArmor string, algo hd.PubKeyType, address types.AccAddress) Info {
	return &localInfo{
		Name:         name,
		PubKey:       pub,
		PrivKeyArmor: privArmor,
		Algo:         algo,
		Address:      address,
	}
}

// GetType implements Info interface
func (i localInfo) GetType() KeyType {
	return TypeLocal
}

// GetType implements Info interface
func (i localInfo) GetName() string {
	return i.Name
}

// GetType implements Info interface
func (i localInfo) GetPubKey() crypto.PubKey {
	return i.PubKey
}

// GetType implements Info interface
func (i localInfo) GetAddress() types.AccAddress {
	return i.Address
}

func (i localInfo) WithNameAndAddress(name string, address types.AccAddress) Info {
	i.Address = address
	i.Name = name
	return i
}

// GetType implements Info interface
func (i localInfo) GetAlgo() hd.PubKeyType {
	return i.Algo
}

// GetType implements Info interface
func (i localInfo) GetPath() (*hd.BIP44Params, error) {
	return nil, fmt.Errorf("BIP44 Paths are not available for this type")
}

// ledgerInfo is the public information about a Ledger key
// Note: Algo must be last field in struct for backwards amino compatibility
type ledgerInfo struct {
	Name    string           `json:"name"`
	PubKey  crypto.PubKey    `json:"pubkey"`
	Path    hd.BIP44Params   `json:"path"`
	Algo    hd.PubKeyType    `json:"algo"`
	Address types.AccAddress `json:"address"`
}

func newLedgerInfo(name string, pub crypto.PubKey, path hd.BIP44Params, algo hd.PubKeyType, address types.AccAddress) Info {
	return &ledgerInfo{
		Name:    name,
		PubKey:  pub,
		Path:    path,
		Algo:    algo,
		Address: address,
	}
}

// GetType implements Info interface
func (i ledgerInfo) GetType() KeyType {
	return TypeLedger
}

// GetName implements Info interface
func (i ledgerInfo) GetName() string {
	return i.Name
}

// GetPubKey implements Info interface
func (i ledgerInfo) GetPubKey() crypto.PubKey {
	return i.PubKey
}

// GetAddress implements Info interface
func (i ledgerInfo) GetAddress() types.AccAddress {
	return i.Address
}

func (i ledgerInfo) WithNameAndAddress(name string, address types.AccAddress) Info {
	i.Address = address
	i.Name = name
	return i
}

// GetPath implements Info interface
func (i ledgerInfo) GetAlgo() hd.PubKeyType {
	return i.Algo
}

// GetPath implements Info interface
func (i ledgerInfo) GetPath() (*hd.BIP44Params, error) {
	tmp := i.Path
	return &tmp, nil
}

// offlineInfo is the public information about an offline key
// Note: Algo must be last field in struct for backwards amino compatibility
type offlineInfo struct {
	Name    string           `json:"name"`
	PubKey  crypto.PubKey    `json:"pubkey"`
	Algo    hd.PubKeyType    `json:"algo"`
	Address types.AccAddress `json:"address"`
}

func newOfflineInfo(name string, pub crypto.PubKey, algo hd.PubKeyType, address types.AccAddress) Info {
	return &offlineInfo{
		Name:    name,
		PubKey:  pub,
		Algo:    algo,
		Address: address,
	}
}

// GetType implements Info interface
func (i offlineInfo) GetType() KeyType {
	return TypeOffline
}

// GetName implements Info interface
func (i offlineInfo) GetName() string {
	return i.Name
}

// GetPubKey implements Info interface
func (i offlineInfo) GetPubKey() crypto.PubKey {
	return i.PubKey
}

// GetAlgo returns the signing algorithm for the key
func (i offlineInfo) GetAlgo() hd.PubKeyType {
	return i.Algo
}

// GetAddress implements Info interface
func (i offlineInfo) GetAddress() types.AccAddress {
	return i.Address
}

func (i offlineInfo) WithNameAndAddress(name string, address types.AccAddress) Info {
	i.Address = address
	i.Name = name
	return i
}

// GetPath implements Info interface
func (i offlineInfo) GetPath() (*hd.BIP44Params, error) {
	return nil, fmt.Errorf("BIP44 Paths are not available for this type")
}

type multisigPubKeyInfo struct {
	PubKey crypto.PubKey `json:"pubkey"`
	Weight uint          `json:"weight"`
}

// multiInfo is the public information about a multisig key
type multiInfo struct {
	Name      string               `json:"name"`
	PubKey    crypto.PubKey        `json:"pubkey"`
	Address   types.AccAddress     `json:"address"`
	Threshold uint                 `json:"threshold"`
	PubKeys   []multisigPubKeyInfo `json:"pubkeys"`
}

// NewMultiInfo creates a new multiInfo instance
func NewMultiInfo(name string, pub crypto.PubKey, address types.AccAddress) Info {
	multiPK := pub.(*multisig.LegacyAminoPubKey)

	pubKeys := make([]multisigPubKeyInfo, len(multiPK.PubKeys))
	for i, pk := range multiPK.GetPubKeys() {
		// TODO: Recursively check pk for total weight?
		pubKeys[i] = multisigPubKeyInfo{pk, 1}
	}

	return &multiInfo{
		Name:      name,
		PubKey:    pub,
		Address:   address,
		Threshold: uint(multiPK.Threshold),
		PubKeys:   pubKeys,
	}
}

// GetType implements Info interface
func (i multiInfo) GetType() KeyType {
	return TypeMulti
}

// GetName implements Info interface
func (i multiInfo) GetName() string {
	return i.Name
}

// GetPubKey implements Info interface
func (i multiInfo) GetPubKey() crypto.PubKey {
	return i.PubKey
}

// GetAddress implements Info interface
func (i multiInfo) GetAddress() types.AccAddress {
	return i.Address
}

func (i multiInfo) WithNameAndAddress(name string, address types.AccAddress) Info {
	i.Address = address
	i.Name = name
	return i
}

// GetPath implements Info interface
func (i multiInfo) GetAlgo() hd.PubKeyType {
	return hd.MultiType
}

// GetPath implements Info interface
func (i multiInfo) GetPath() (*hd.BIP44Params, error) {
	return nil, fmt.Errorf("BIP44 Paths are not available for this type")
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (i multiInfo) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	multiPK := i.PubKey.(*multisig.LegacyAminoPubKey)

	return codectypes.UnpackInterfaces(multiPK, unpacker)
}

// encoding info
func marshalInfo(i Info) []byte {
	return CryptoCdc.MustMarshalBinaryLengthPrefixed(i)
}

// decoding info
func unmarshalInfo(bz []byte) (info Info, err error) {
	err = CryptoCdc.UnmarshalBinaryLengthPrefixed(bz, &info)
	if err != nil {
		return nil, err
	}

	// After unmarshalling into &info, if we notice that the info is a
	// multiInfo, then we unmarshal again, explicitly in a multiInfo this time.
	// Since multiInfo implements UnpackInterfacesMessage, this will correctly
	// unpack the underlying anys inside the multiInfo.
	//
	// This is a workaround, as go cannot check that an interface (Info)
	// implements another interface (UnpackInterfacesMessage).
	_, ok := info.(multiInfo)
	if ok {
		var multi multiInfo
		err = CryptoCdc.UnmarshalBinaryLengthPrefixed(bz, &multi)

		return multi, err
	}

	return
}
