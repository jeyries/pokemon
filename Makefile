
ALL := \
 pokemon_swift.exe \
 pokemon_c.exe \
 pokemon_go.exe \
 pokemon_rust.exe

all: $(ALL)

clean:
	rm -f $(ALL)

###### Swift

SWIFT_FLAGS := -Ounchecked -Xcc -march=native -Xllvm -slp-threshold=1000000 -Xllvm -unroll-count=1

pokemon_swift.asm: pokemon.swift
	swiftc -emit-assembly $(SWIFT_FLAGS) $< -o $@
	xcrun swift-demangle < $@ > pokemon.demangle.asm

pokemon_swift.exe: pokemon.swift
	swiftc -emit-executable $(SWIFT_FLAGS) $< -o $@

###### C

#CFLAGS := -O3 -march=native
#CC := gcc-8
CFLAGS := -O3 -march=native -fno-slp-vectorize
CC := clang

pokemon_c.asm: pokemon.c
	$(CC) -S $(CFLAGS) $< -o $@

pokemon_c.exe: pokemon.c
	$(CC) $(CFLAGS) $< -o $@

###### Go

pokemon_go.s: pokemon.go
	go tool compile -S $< > $@

pokemon_go.exe: pokemon.go
	go build -o $@ $<

###### Rust

RUSTFLAGS := -O -C target-cpu=native

pokemon_rust.asm: pokemon.rs
	rustc pokemon.rs $(RUSTFLAGS) --emit=asm -o pokemon_rust.asm
	
pokemon_rust.exe: pokemon.rs
	rustc pokemon.rs $(RUSTFLAGS) -o pokemon_rust.exe



