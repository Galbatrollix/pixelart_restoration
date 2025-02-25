package common

func Ternary[T any](condition bool, ifTrue T, ifFalse T) T {
    if condition {
        return ifTrue
    }else{
        return ifFalse
    }
    
}