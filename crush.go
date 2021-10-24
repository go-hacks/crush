//File hasher for when speed is more important
//than security and we just want to verify
//that file transfers were complete.
//Approx. 3X+ faster than md5sum w/ half the bits.
//Only non-core hash values escape to heap.

package main

import (
  "bufio"
  "encoding/hex"
  "encoding/binary"
  "fmt"
  "os"
)

//Basic world values
const readSize int = 4096000
const stateByteCnt = 8

func main () {
  //Input argument sanity check
  if len(os.Args) != 2 {
    fmt.Println("Usage: crush filename")
    os.Exit(0)
  }

  //File sanity checking
  fileCheck, _ := os.Stat(os.Args[1])
  if fileCheck == nil {
    os.Stderr.WriteString("File does not exist!\n")
    os.Exit(0)
  } else if fileCheck.IsDir() {
    os.Stderr.WriteString("Given file is a directory!\n")
    os.Exit(0)
  }

  //Open file for reading
  inFile, _ := os.OpenFile(os.Args[1],os.O_RDONLY,0666)
  defer inFile.Close()

  //Make read buffer
  bufRdr := bufio.NewReaderSize(inFile, readSize)

  //Initial state (sqrt(3))
  stateHex := "1c98c677de371c7d"
  stateBytes := make([]byte, stateByteCnt)
	hex.Decode(stateBytes, []byte(stateHex))
  state := binary.BigEndian.Uint64(stateBytes)

  //Buffer for reads into hasher
  data := make([]byte, readSize)

  //Length for handling block size
  length := readSize

  //Hashing core XORs every 64bits & fills incomplete blocks
  for {
    n, _ := bufRdr.Read(data)
    if n == 0 {
      break
    }
    //If not 64bits fill with 01010101
    if n % stateByteCnt != 0 {
      fillCnt := stateByteCnt - (n % stateByteCnt)
      length = n + fillCnt
      for i := 0; i < fillCnt; i++ {
        data[n+i] = byte(55)
      }
    }
    //XOR every 64bits with the state
    for i := 0; i < length / stateByteCnt; i++ {
      val64 := binary.BigEndian.Uint64(data[i*stateByteCnt:i*stateByteCnt+stateByteCnt])
      state ^= val64
    }
  }

  //Convert 64bit uint back to bytes for nlfsr
  binary.BigEndian.PutUint64(stateBytes[0:stateByteCnt], state)

  //Extrapolate state by use of nlfsr
  nlfsr(stateBytes)

  //Convert raw bytes back into hex string for printing
  hashHexBytes := make([]byte, hex.EncodedLen(len(stateBytes)))
	hex.Encode(hashHexBytes, stateBytes)
  hashStr := fmt.Sprintf("%s", hashHexBytes)

  //Print output
  if hashStr == "b4bb4023dcbf444b" {
    os.Stderr.WriteString("Hash is the base state because the file\n")
    os.Stderr.WriteString("is all zeroes and divisible by 64bits!\n")
    fmt.Printf("%s %s\n", hashHexBytes, os.Args[1])
  } else {
    fmt.Printf("%s %s\n", hashHexBytes, os.Args[1])
  }
  return
}

//Non-Linear Feed Shift Register for extrapolating raw state.
//This makes every input bit affect the rest of the state.
func nlfsr (src []byte) []byte {
  dst := make([]byte, stateByteCnt)
  for i := 0; i < stateByteCnt; i++ {
    for j := 0; j < stateByteCnt; j++ {
      dst = append(src[1:8],src[0]^src[1]^src[2]^src[3]^src[4]^src[5]^src[6]^src[7])
      copy(src,dst)
      src[0] += dst[7]
      src[1] += dst[6]
      src[2] += dst[5]
      src[3] += dst[4]
      src[4] += dst[3]
      src[5] += dst[2]
      src[6] += dst[1]
      src[7] += dst[0]
    }
    src[0] ^= 255
    dst[1] ^= 255
    src[2] ^= 255
    dst[3] ^= 255
    src[4] ^= 255
    dst[5] ^= 255
    src[6] ^= 255
    dst[7] ^= 255
  }
  return dst
}
