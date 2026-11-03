import HikeCard from "../../../pages/HikeLog/HikeCard/HikeCard";
import { HIKE_MOCK_1 } from "../../mocks/models/hike";

/**
 * Wraps around a HikeCard to provide default prop values for testing.
 */
export function HikeCardWrapper({ isExpanded = false }: { isExpanded?: boolean }) {
  return (
    <HikeCard
      deleteEnabled={true}
      hike={HIKE_MOCK_1}
      index={1}
      isExpanded={isExpanded}
      onClick={() => {}}
      onDelete={async () => {}}
    />
  );
}
