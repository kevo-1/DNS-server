package protocol

type Message struct {
    Header      Header
    Questions   []Question
    Answers     []ResourceRecord
    Authorities []ResourceRecord
    Additional  []ResourceRecord
}

type Header struct {
    ID              uint16
    Flags           uint16
    QuestionCount   uint16
    AnswerCount     uint16
    AuthorityCount  uint16
    AdditionalCount uint16
}

type Question struct {
    Name  string
    Type  uint16
    Class uint16
}

type ResourceRecord struct {
    Name     string
    Type     uint16
    Class    uint16
    TTL      uint32
    RDLength uint16
    RData    []byte
}