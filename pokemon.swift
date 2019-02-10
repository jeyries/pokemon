//
//  main.swift
//  Pokemon
//
//  Created by Julien on 21/09/2018.
//  Copyright © 2018 Julien Eyriès. All rights reserved.
//

import Foundation

import simd

/// Read text file line by line
public class LineReader {
    public let path: String

    fileprivate let file: UnsafeMutablePointer<FILE>!

    init?(path: String) {
        self.path = path
        file = fopen(path, "r")
        guard file != nil else { return nil }
    }

    public var nextLine: String? {
        var line:UnsafeMutablePointer<CChar>? = nil
        var linecap:Int = 0
        defer { free(line) }
        return getline(&line, &linecap, file) > 0 ? String(cString: line!) : nil
    }

    deinit {
        fclose(file)
    }
}

extension LineReader: Sequence {
    public func  makeIterator() -> AnyIterator<String> {
        return AnyIterator<String> {
            return self.nextLine
        }
    }
}

///////////////

struct Pokemon {

    let name: String
    let attacks: [UInt8]

    func setbits( bits: inout [UInt64], offset: Int ) {
        for i in 0..<3 {
            bits[offset + i] = 0
        }
        for a in attacks {
            bits[offset + Int(a/64)] |= 1 << (a % 64)
        }
    }


}

@inline(never)
func search( bits: UnsafeBufferPointer<UInt64>, m: Int, answer: UnsafeMutableBufferPointer<Int> )  {

    var best = 0

    for _i0 in 0..<m {
        let i0 = m-1 - _i0
        print(String(format: "processing i0= %d", i0));
        let start = clock();

        for i1 in 0..<i0 {
            for i2 in 0..<i1 {
                for i3 in 0..<i2 {
                    for i4 in 0..<i3 {

                        var s0: UInt64 = 0
                        var s1: UInt64 = 0
                        var s2: UInt64 = 0
                        s0 |= bits[3*i0 + 0]
                        s1 |= bits[3*i0 + 1]
                        s2 |= bits[3*i0 + 2]

                        s0 |= bits[3*i1 + 0]
                        s1 |= bits[3*i1 + 1]
                        s2 |= bits[3*i1 + 2]

                        s0 |= bits[3*i2 + 0]
                        s1 |= bits[3*i2 + 1]
                        s2 |= bits[3*i2 + 2]

                        s0 |= bits[3*i3 + 0]
                        s1 |= bits[3*i3 + 1]
                        s2 |= bits[3*i3 + 2]

                        s0 |= bits[3*i4 + 0]
                        s1 |= bits[3*i4 + 1]
                        s2 |= bits[3*i4 + 2]
                        
                        for i5 in 0..<i4 {
                            //print("\(i0) \(i1) \(i2) ")

                            let _s0 = s0 | bits[3*i5 + 0]
                            let _s1 = s1 | bits[3*i5 + 1]
                            let _s2 = s2 | bits[3*i5 + 2]

                            let score = _s0.nonzeroBitCount + _s1.nonzeroBitCount + _s2.nonzeroBitCount
                            if score > best {
                                best = score
                                answer[0] = i0
                                answer[1] = i1
                                answer[2] = i2
                                answer[3] = i3
                                answer[4] = i4
                                answer[5] = i5
                            }
                        }
                    }
                }
            }
        }

        let end = clock()
        let elapsed = Double( end - start ) / Double(CLOCKS_PER_SEC)
        print(String(format: "elapsed = %.3f s", elapsed))
        print(String(format: "best = %d", best))
    }

}

func main() {

    assert( CommandLine.arguments.count == 2 )
    let path = CommandLine.arguments[1]

    guard let reader = LineReader(path: path) else {
        preconditionFailure("cannot open file \(path)");
    }

    var pokemons = [Pokemon]()
    for line in reader {
        let items = line.trimmingCharacters(in: .whitespacesAndNewlines).split(separator: ",")
        let name = String(items[0])
        let attacks = items[1...].map { UInt8($0)! }
        let pokemon = Pokemon(name: name, attacks: attacks)
        pokemons.append(pokemon)
    }

    print(pokemons[0])

    var bits = [UInt64](repeating: 0, count: 3*pokemons.count)
    for i in pokemons.indices {
        pokemons[i].setbits( bits: &bits, offset: i*3)
    }

    print(String(format:"%x", bits[0]))

    var answer = [Int](repeating: 0, count: 6)
    answer.withUnsafeMutableBufferPointer { answer in
        bits.withUnsafeBufferPointer { bits in
            search( bits: bits, m: pokemons.count, answer: answer)
        }
    }


    for a in answer {
        print(pokemons[a].name)
    }


}

///

main()
