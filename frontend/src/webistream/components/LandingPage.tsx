import {
  Box,
  Button,
  Container,
  Flex,
  Heading,
  Image,
  Text,
} from "@chakra-ui/react"
import {
  ArrowRight,
  Briefcase,
  Calendar,
  LayoutDashboard,
  LogOut,
  PlayCircle,
  Video,
  Zap,
} from "lucide-react"
import type React from "react"
import { EventStatus, type User, UserRole, type WebinarEvent } from "../types"

interface LandingPageProps {
  events: WebinarEvent[]
  currentUser: User | null
  onHostLogin: () => void
  onMyWebinars: () => void
  onRegister: (event: WebinarEvent) => void
  onAuthClick: (role?: UserRole) => void
  onLogout: () => void
}

const LandingPage: React.FC<LandingPageProps> = ({
  events,
  currentUser,
  onHostLogin,
  onMyWebinars,
  onRegister,
  onAuthClick,
  onLogout,
}) => {
  const publishedEvents = events.filter((e) => e.status !== EventStatus.DRAFT)

  const scrollTo = (id: string) => {
    const element = document.getElementById(id)
    if (element) {
      element.scrollIntoView({ behavior: "smooth" })
    }
  }

  return (
    <Box bg="white" color="gray.900" overflowX="hidden">
      <Flex
        as="nav"
        position="fixed"
        top={0}
        w="full"
        zIndex={50}
        bg="whiteAlpha.800"
        backdropFilter="auto"
        backdropBlur="md"
        borderBottom="1px solid"
        borderColor="gray.100"
        px={6}
        py={4}
        align="center"
        justify="space-between"
      >
        <Button variant="ghost" onClick={() => scrollTo("hero")} gap={2} px={2}>
          <Flex
            bg="blue.600"
            p={2}
            borderRadius="lg"
            transform="auto"
            _hover={{ rotate: "12deg" }}
          >
            <Video color="white" size={20} />
          </Flex>
          <Text fontWeight="bold" fontSize="xl" letterSpacing="tight">
            WebiStream
          </Text>
        </Button>
        <Flex
          display={{ base: "none", md: "flex" }}
          align="center"
          gap={8}
          fontSize="sm"
          fontWeight="medium"
          color="gray.600"
        >
          <Button
            variant="ghost"
            onClick={() => scrollTo("about")}
            color="gray.600"
            _hover={{ color: "blue.600" }}
          >
            About
          </Button>
          <Button
            variant="ghost"
            onClick={() => scrollTo("webinars")}
            color="gray.600"
            _hover={{ color: "blue.600" }}
          >
            Webinars
          </Button>
          <Button
            variant="ghost"
            onClick={() => scrollTo("business")}
            color="gray.900"
            fontWeight="bold"
            _hover={{ color: "blue.600" }}
          >
            For Businesses
          </Button>
        </Flex>
        <Flex align="center" gap={4}>
          {!currentUser ? (
            <>
              <Button
                variant="ghost"
                onClick={() => onAuthClick()}
                fontSize="sm"
                fontWeight="semibold"
                color="gray.700"
                _hover={{ color: "blue.600" }}
              >
                Sign In
              </Button>
              <Button
                onClick={onHostLogin}
                colorPalette="gray"
                variant="solid"
                borderRadius="full"
                fontSize="sm"
                fontWeight="bold"
              >
                Host Portal
              </Button>
            </>
          ) : (
            <Flex align="center" gap={4}>
              {currentUser.role === UserRole.BUSINESS ? (
                <Button
                  variant="ghost"
                  onClick={onHostLogin}
                  colorPalette="blue"
                  fontWeight="bold"
                  px={4}
                  borderRadius="xl"
                >
                  <LayoutDashboard size={16} />
                  Go to Dashboard
                </Button>
              ) : (
                <Button
                  variant="ghost"
                  onClick={onMyWebinars}
                  colorPalette="blue"
                  fontWeight="bold"
                  px={4}
                  borderRadius="xl"
                >
                  <PlayCircle size={16} />
                  My Learning
                </Button>
              )}
              <Flex position="relative" role="group">
                <Image
                  src={currentUser.avatar}
                  w={10}
                  h={10}
                  borderRadius="full"
                  border="2px solid"
                  borderColor="blue.100"
                />
                <Box
                  position="absolute"
                  right={0}
                  mt={2}
                  w={48}
                  bg="white"
                  borderRadius="2xl"
                  boxShadow="xl"
                  border="1px solid"
                  borderColor="gray.100"
                  py={2}
                  display="none"
                  _groupHover={{ display: "block" }}
                >
                  <Box
                    px={4}
                    py={2}
                    borderBottom="1px solid"
                    borderColor="gray.50"
                    mb={2}
                  >
                    <Text
                      fontSize="xs"
                      fontWeight="bold"
                      color="gray.400"
                      letterSpacing="widest"
                      textTransform="uppercase"
                    >
                      Signed in as
                    </Text>
                    <Text fontSize="sm" fontWeight="bold" color="gray.900">
                      {currentUser.name}
                    </Text>
                  </Box>
                  <Button
                    onClick={onLogout}
                    variant="ghost"
                    colorPalette="red"
                    justifyContent="flex-start"
                    w="full"
                    px={4}
                    py={2}
                    fontSize="sm"
                    fontWeight="bold"
                    gap={2}
                  >
                    <LogOut size={16} />
                    Sign Out
                  </Button>
                </Box>
              </Flex>
            </Flex>
          )}
        </Flex>
      </Flex>

      <Box
        id="hero"
        pt={32}
        pb={20}
        px={6}
        bgGradient="to-b"
        gradientFrom="blue.50"
        gradientTo="white"
      >
        <Container maxW="6xl">
          <Flex direction="column" align="center" textAlign="center" gap={8}>
            <Flex
              align="center"
              gap={2}
              px={3}
              py={1}
              bg="blue.100"
              color="blue.700"
              borderRadius="full"
              fontSize="xs"
              fontWeight="bold"
              textTransform="uppercase"
              letterSpacing="widest"
              w="fit-content"
            >
              <Zap size={12} /> The Future of Knowledge Sharing
            </Flex>
            <Heading
              fontSize={{ base: "5xl", md: "7xl" }}
              fontWeight="extrabold"
              color="gray.900"
            >
              Learn from the{" "}
              <Text as="span" color="blue.600">
                World's Best
              </Text>{" "}
              <br /> Minds, Live.
            </Heading>
            <Text maxW="2xl" mx="auto" fontSize="xl" color="gray.500">
              WebiStream connects industry experts with eager learners through
              high-fidelity live streaming, interactive workshops, and secure
              payment processing.
            </Text>
            <Flex
              direction={{ base: "column", sm: "row" }}
              align="center"
              justify="center"
              gap={4}
              pt={4}
              w="full"
            >
              <Button
                onClick={() => scrollTo("webinars")}
                variant="solid"
                colorPalette="gray"
                px={8}
                py={4}
                borderRadius="2xl"
                fontWeight="bold"
                boxShadow="xl"
                gap={2}
              >
                Browse All Events <ArrowRight size={20} />
              </Button>
              <Button
                onClick={() => onAuthClick(UserRole.BUSINESS)}
                variant="outline"
                colorPalette="blue"
                px={8}
                py={4}
                borderRadius="2xl"
                fontWeight="bold"
                gap={2}
              >
                <Briefcase size={20} /> Host a Webinar
              </Button>
            </Flex>
          </Flex>
        </Container>
      </Box>

      <Box id="webinars" py={24} px={6} bg="gray.50">
        <Container maxW="6xl">
          <Flex
            direction={{ base: "column", md: "row" }}
            justify="space-between"
            align="flex-end"
            gap={6}
            mb={12}
          >
            <Box>
              <Heading fontSize="4xl" mb={4} letterSpacing="tight">
                Upcoming Live Sessions
              </Heading>
              <Text color="gray.500" fontSize="lg">
                Don't miss out on these trending topics.
              </Text>
            </Box>
          </Flex>

          <Box
            display="grid"
            gridTemplateColumns={{
              base: "1fr",
              md: "repeat(2, 1fr)",
              lg: "repeat(3, 1fr)",
            }}
            gap={8}
          >
            {publishedEvents.map((event) => {
              const isRegistered =
                currentUser &&
                event.attendees.some((a) => a.email === currentUser.email)
              return (
                <Box
                  key={event.id}
                  bg="white"
                  borderRadius="3xl"
                  overflow="hidden"
                  border="1px solid"
                  borderColor="gray.100"
                  boxShadow="sm"
                  _hover={{ boxShadow: "2xl" }}
                  transition="all"
                  display="flex"
                  flexDirection="column"
                >
                  <Box
                    position="relative"
                    h={48}
                    bg="gray.900"
                    overflow="hidden"
                  >
                    <Image
                      src={`https://picsum.photos/seed/${event.id}/800/600`}
                      w="full"
                      h="full"
                      objectFit="cover"
                      opacity={0.6}
                    />
                    <Box
                      position="absolute"
                      top={4}
                      left={4}
                      px={3}
                      py={1}
                      bg="whiteAlpha.200"
                      backdropFilter="auto"
                      backdropBlur="md"
                      borderRadius="full"
                      color="white"
                      fontSize="xs"
                      fontWeight="black"
                      letterSpacing="widest"
                      textTransform="uppercase"
                    >
                      {event.platform}
                    </Box>
                  </Box>
                  <Flex p={8} gap={4} direction="column" flex="1">
                    <Flex
                      align="center"
                      gap={2}
                      color="blue.600"
                      fontWeight="bold"
                      fontSize="xs"
                      textTransform="uppercase"
                      letterSpacing="widest"
                    >
                      <Calendar size={12} />
                      {new Date(event.date).toLocaleDateString()}
                    </Flex>
                    <Text
                      fontSize="xl"
                      fontWeight="bold"
                      _groupHover={{ color: "blue.600" }}
                      flex="1"
                    >
                      {event.title}
                    </Text>
                    <Flex
                      align="center"
                      justify="space-between"
                      pt={4}
                      borderTop="1px solid"
                      borderColor="gray.50"
                    >
                      <Text fontSize="2xl" fontWeight="black" color="gray.900">
                        ${event.price}
                      </Text>
                      <Text
                        as="span"
                        fontSize="xs"
                        fontWeight="bold"
                        color="gray.400"
                        letterSpacing="widest"
                        textTransform="uppercase"
                      >
                        {event.attendees.length} enrolled
                      </Text>
                    </Flex>
                    <Button
                      onClick={() => onRegister(event)}
                      disabled={
                        !!isRegistered ||
                        currentUser?.role === UserRole.BUSINESS
                      }
                      borderRadius="2xl"
                      fontWeight="bold"
                      variant={isRegistered ? "subtle" : "solid"}
                      colorPalette={isRegistered ? "green" : "blue"}
                    >
                      {isRegistered
                        ? "Already Enrolled"
                        : currentUser?.role === UserRole.BUSINESS
                          ? "Host Account"
                          : "Register Now"}
                    </Button>
                  </Flex>
                </Box>
              )
            })}
          </Box>
        </Container>
      </Box>

      <Box id="business" py={24} px={6}>
        <Container maxW="6xl">
          <Box
            bg="gray.900"
            borderRadius="3rem"
            p={{ base: 12, md: 20 }}
            color="white"
            overflow="hidden"
            position="relative"
          >
            <Box position="relative" zIndex={10} maxW="2xl">
              <Box
                px={4}
                py={1.5}
                bg="blue.600"
                w="fit"
                borderRadius="full"
                fontSize="xs"
                fontWeight="black"
                textTransform="uppercase"
                letterSpacing="wide"
              >
                For Businesses & Creators
              </Box>
              <Heading
                fontSize={{ base: "4xl", md: "6xl" }}
                fontWeight="black"
                lineHeight="tight"
                mt={6}
              >
                Scale your knowledge into a revenue stream.
              </Heading>
              <Text fontSize="xl" color="gray.400" mt={4}>
                Connect your Stripe account, sync with Zoom or Google Meet, and
                start hosting paid webinars in minutes. Our AI handles the
                descriptions and pricing suggestions.
              </Text>
              <Flex direction={{ base: "column", sm: "row" }} gap={4} mt={8}>
                <Button
                  onClick={() => onAuthClick(UserRole.BUSINESS)}
                  variant="solid"
                  colorPalette="blue"
                  px={10}
                  py={5}
                  borderRadius="2xl"
                  fontWeight="black"
                  boxShadow="2xl"
                  gap={3}
                >
                  Start Hosting Now <ArrowRight size={24} />
                </Button>
              </Flex>
            </Box>
          </Box>
        </Container>
      </Box>
    </Box>
  )
}

export default LandingPage
