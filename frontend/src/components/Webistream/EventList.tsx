import {
  EventStatus,
  PlatformType,
  type WebinarEvent,
} from "@/types/webistream"
import {
  Box,
  Button,
  Flex,
  Grid,
  Heading,
  IconButton,
  Text,
} from "@chakra-ui/react"
import {
  Calendar as CalendarIcon,
  Globe,
  Layers,
  Lock,
  MoreVertical,
  Users,
  Video,
  Zap,
} from "lucide-react"
import type React from "react"

interface EventListProps {
  events: WebinarEvent[]
  onManageParticipants: (event: WebinarEvent) => void
  onPublish: (event: WebinarEvent) => void
  onOpenControlPanel: (event: WebinarEvent) => void
}

const EventList: React.FC<EventListProps> = ({
  events,
  onManageParticipants,
  onPublish,
  onOpenControlPanel,
}) => {
  return (
    <Grid
      gridTemplateColumns={{
        base: "1fr",
        md: "repeat(2, 1fr)",
        xl: "repeat(3, 1fr)",
      }}
      gap={8}
    >
      {events.map((event) => (
        <EventCard
          key={event.id}
          event={event}
          onManageParticipants={() => onManageParticipants(event)}
          onPublish={() => onPublish(event)}
          onControlPanel={() => onOpenControlPanel(event)}
        />
      ))}
    </Grid>
  )
}

const EventCard: React.FC<{
  event: WebinarEvent
  onManageParticipants: () => void
  onPublish: () => void
  onControlPanel: () => void
}> = ({ event, onManageParticipants, onPublish, onControlPanel }) => {
  const isDraft = event.status === EventStatus.DRAFT
  const headerBg = isDraft
    ? "amber.500"
    : event.platform === PlatformType.ZOOM
      ? "indigo.600"
      : "emerald.600"
  const statusColor = isDraft ? "amber.600" : "indigo.600"

  return (
    <Box
      bg="white"
      borderRadius="2rem"
      border="1px solid"
      borderColor={isDraft ? "amber.100" : "gray.100"}
      transition="all 0.3s ease"
      overflow="hidden"
      display="flex"
      flexDirection="column"
      _hover={{
        borderColor: isDraft ? "amber.300" : "indigo.300",
        boxShadow: isDraft ? "sm" : "2xl",
        transform: isDraft ? undefined : "translateY(-4px)",
      }}
    >
      <Box
        position="relative"
        h={40}
        display="flex"
        flexDirection="column"
        p={6}
        justifyContent="space-between"
        bg={headerBg}
      >
        <Flex justify="space-between" align="flex-start">
          <Box
            bg="whiteAlpha.200"
            backdropFilter="blur(8px)"
            p={2.5}
            borderRadius="2xl"
            color="white"
          >
            <Video size={24} />
          </Box>
          <Flex direction="column" align="flex-end" gap={2}>
            <Text
              px={3}
              py={1}
              bg="whiteAlpha.200"
              backdropFilter="blur(8px)"
              color="white"
              fontSize="10px"
              fontWeight="black"
              textTransform="uppercase"
              letterSpacing="widest"
              borderRadius="full"
            >
              {event.platform}
            </Text>
            <Text
              px={3}
              py={1}
              bg="white"
              fontSize="10px"
              fontWeight="black"
              textTransform="uppercase"
              letterSpacing="widest"
              borderRadius="full"
              boxShadow="sm"
              color={statusColor}
            >
              {event.status}
            </Text>
          </Flex>
        </Flex>

        <Flex align="center" gap={2}>
          {event.meetingLink ? (
            <Flex
              align="center"
              gap={1.5}
              px={3}
              py={1}
              bg="whiteAlpha.200"
              backdropFilter="blur(8px)"
              color="white"
              borderRadius="full"
              fontSize="xs"
              fontWeight="medium"
            >
              <Globe size={12} /> Link Ready
            </Flex>
          ) : (
            <Flex
              align="center"
              gap={1.5}
              px={3}
              py={1}
              bg="blackAlpha.300"
              backdropFilter="blur(8px)"
              color="whiteAlpha.700"
              borderRadius="full"
              fontSize="xs"
              fontWeight="medium"
            >
              <Lock size={12} /> Setup Needed
            </Flex>
          )}
        </Flex>
      </Box>

      <Box p={8} flex="1" display="flex" flexDirection="column">
        <Flex justify="space-between" align="flex-start" mb={4}>
          <Heading size="md" color="gray.900">
            {event.title}
          </Heading>
          <IconButton
            aria-label="Options"
            variant="ghost"
            size="sm"
            color="gray.400"
            _hover={{ color: "gray.600" }}
          >
            <MoreVertical size={20} />
          </IconButton>
        </Flex>

        <Text color="gray.500" fontSize="sm" mb={6}>
          {event.description}
        </Text>

        <Grid
          gridTemplateColumns="repeat(2, 1fr)"
          gap={4}
          mb={8}
          fontSize="xs"
          fontWeight="bold"
          color="gray.400"
          textTransform="uppercase"
          letterSpacing="widest"
        >
          <Flex align="center" gap={2}>
            <CalendarIcon size={16} color="#818cf8" />
            <Text color="gray.600">
              {new Date(event.date).toLocaleDateString()}
            </Text>
          </Flex>
          <Flex align="center" gap={2}>
            <Users size={16} color="#818cf8" />
            <Text color="gray.600">{event.attendees.length} Users</Text>
          </Flex>
        </Grid>

        <Flex
          pt={6}
          borderTop="1px solid"
          borderColor="gray.100"
          align="center"
          justify="space-between"
          gap={4}
        >
          {isDraft ? (
            <Button
              onClick={onPublish}
              flex="1"
              px={6}
              py={3}
              borderRadius="2xl"
              fontSize="sm"
              fontWeight="black"
              colorPalette="amber"
              variant="solid"
              boxShadow="lg"
              display="inline-flex"
              alignItems="center"
              gap={2}
              justifyContent="center"
            >
              <Zap size={16} /> Publish Now ($10)
            </Button>
          ) : (
            <>
              <Button
                onClick={onManageParticipants}
                flex="1"
                px={4}
                py={3}
                borderRadius="2xl"
                fontSize="sm"
                fontWeight="black"
                variant="subtle"
                colorPalette="gray"
                display="inline-flex"
                alignItems="center"
                gap={2}
                justifyContent="center"
              >
                <Users size={16} /> Manage
              </Button>
              <IconButton
                aria-label="Open Control Panel"
                onClick={onControlPanel}
                borderRadius="2xl"
                colorPalette="indigo"
                variant="solid"
                title="Open Control Panel"
              >
                <Layers size={20} />
              </IconButton>
            </>
          )}
        </Flex>
      </Box>
    </Box>
  )
}

export default EventList
