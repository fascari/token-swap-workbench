import { useCallback, useRef, useState } from "react";

import type { Block } from "../api/types";
import { fetchBlocks } from "../api/client";

export function useBlocks() {
  const [blocks, setBlocks] = useState<Block[] | null>(null);
  const [blockCount, setBlockCount] = useState(10);
  const countRef = useRef(blockCount);
  countRef.current = blockCount;

  const load = useCallback(async (count?: number) => {
    const n = count ?? countRef.current;
    const result = await fetchBlocks(n);
    setBlocks(result);
    return result;
  }, []);

  return { blocks, blockCount, setBlockCount, load };
}
