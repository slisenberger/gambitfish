# Gambitfish Chess Engine

## What

Gambitfish is a work-in-progress chess engine designed by Stefan Isenberger. It is incomplete in both correctness and evaluation.

## Features

Gambitfish is currently not particularly robust; its evaluator is simple material calculation and search is basic alpha beta. However, ultimately the goal is to have a robust chess engine with fancy search.

Gambitfish's "character" will be from preferring early-game gambits to "solid" play. This will weaken it compared to other engines, but hopefully it will be capable to compensate against stronger human players.

Expected implementation of the gambit preference is likely to be an opening-book with gambit preferred lines, however it is possible that pawn material disadvantages may be discounted in the early game.
