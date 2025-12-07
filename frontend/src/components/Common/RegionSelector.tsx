import { Button, Listbox, Popover, Portal, createListCollection } from "@chakra-ui/react"
import { LuChevronDown } from "react-icons/lu"
import { useState, useRef } from "react"

// Common regions with their display names
const REGIONS = [
  { value: "ID", label: "Indonesia (+62)" },
  { value: "US", label: "United States (+1)" },
  { value: "GB", label: "United Kingdom (+44)" },
  { value: "SG", label: "Singapore (+65)" },
  { value: "MY", label: "Malaysia (+60)" },
  { value: "TH", label: "Thailand (+66)" },
  { value: "PH", label: "Philippines (+63)" },
  { value: "VN", label: "Vietnam (+84)" },
  { value: "AU", label: "Australia (+61)" },
  { value: "IN", label: "India (+91)" },
  { value: "CN", label: "China (+86)" },
  { value: "JP", label: "Japan (+81)" },
  { value: "KR", label: "South Korea (+82)" },
  { value: "DE", label: "Germany (+49)" },
  { value: "FR", label: "France (+33)" },
  { value: "IT", label: "Italy (+39)" },
  { value: "ES", label: "Spain (+34)" },
  { value: "NL", label: "Netherlands (+31)" },
  { value: "BE", label: "Belgium (+32)" },
  { value: "CH", label: "Switzerland (+41)" },
  { value: "AT", label: "Austria (+43)" },
  { value: "SE", label: "Sweden (+46)" },
  { value: "NO", label: "Norway (+47)" },
  { value: "DK", label: "Denmark (+45)" },
  { value: "FI", label: "Finland (+358)" },
  { value: "PL", label: "Poland (+48)" },
  { value: "BR", label: "Brazil (+55)" },
  { value: "MX", label: "Mexico (+52)" },
  { value: "CA", label: "Canada (+1)" },
  { value: "NZ", label: "New Zealand (+64)" },
]

interface RegionSelectorProps {
  value?: string
  onChange: (value: string) => void
  disabled?: boolean
  placeholder?: string
}

export const RegionSelector = ({ value, onChange, disabled, placeholder = "Select Region" }: RegionSelectorProps) => {
  const [open, setOpen] = useState(false)
  const triggerRef = useRef<HTMLButtonElement | null>(null)

  const collection = createListCollection({
    items: REGIONS,
  })

  const selectedRegion = REGIONS.find((r) => r.value === value)

  return (
    <Listbox.Root
      collection={collection}
      value={value ? [value] : []}
      onValueChange={(details) => {
        if (details.value && details.value.length > 0) {
          onChange(details.value[0] as string)
        }
        setOpen(false)
      }}
      disabled={disabled}
    >
      <Popover.Root open={open} onOpenChange={(e) => setOpen(e.open)}>
        <Popover.Trigger asChild>
          <Button
            ref={triggerRef}
            variant="outline"
            size="md"
            w="100%"
            justifyContent="space-between"
            disabled={disabled}
          >
            <span>{selectedRegion ? selectedRegion.label : placeholder}</span>
            <LuChevronDown />
          </Button>
        </Popover.Trigger>
        <Portal>
          <Popover.Positioner>
            <Popover.Content maxH="300px" overflowY="auto">
              <Popover.Body p="0">
                <Listbox.Content>
                  {collection.items.map((item) => (
                    <Listbox.Item key={item.value as string} item={item}>
                      <Listbox.ItemText>{item.label as string}</Listbox.ItemText>
                      <Listbox.ItemIndicator />
                    </Listbox.Item>
                  ))}
                </Listbox.Content>
              </Popover.Body>
            </Popover.Content>
          </Popover.Positioner>
        </Portal>
      </Popover.Root>
    </Listbox.Root>
  )
}

