import { useState } from "react";

import HikeCard from "../../../pages/HikeLog/HikeCard";
import { HIKE_MOCK_1 } from "../../mocks/models/hike";

/**
 * Wrapper around a HikeCard which allows it to be expanded/unexpanded.
 */
export function StatefulHikeCard() {
  const [isExpanded, setIsExpanded] = useState<boolean>(false);

  const handleClick = () => {
    setIsExpanded(!isExpanded);
  };

  return (
    <HikeCard hike={HIKE_MOCK_1} index={1} isExpanded={isExpanded} onClick={handleClick} onDelete={async () => {}} />
  );
}

/**
 * Wraps around a HikeCard to provide default prop values for testing.
 */
export function HikeCardWrapper({ isExpanded = false }: { isExpanded?: boolean }) {
  return <HikeCard hike={HIKE_MOCK_1} index={1} isExpanded={isExpanded} onClick={() => {}} onDelete={async () => {}} />;
}
