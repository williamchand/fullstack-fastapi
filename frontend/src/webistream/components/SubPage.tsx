import { Box, Button, Container, Flex, Heading, Text } from "@chakra-ui/react"
import { ArrowLeft, ArrowRight, Check, ExternalLink, Video } from "lucide-react"
import type React from "react"

interface SubPageProps {
  pageKey: string
  onBack: () => void
  onHostLogin: () => void
}

const SubPage: React.FC<SubPageProps> = ({ pageKey, onBack, onHostLogin }) => {
  // removed unused flag

  const renderContent = () => {
    switch (pageKey) {
      case "Pricing":
        return (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <PriceCard
              tier="Starter"
              price="Free"
              items={[
                "Up to 10 attendees",
                "Standard quality",
                "Chat features",
              ]}
            />
            <PriceCard
              tier="Pro"
              price="$29"
              highlight
              items={[
                "Up to 100 attendees",
                "HD Streaming",
                "Recording features",
                "AI Assistant",
              ]}
            />
            <PriceCard
              tier="Enterprise"
              price="Custom"
              items={[
                "Unlimited attendees",
                "Dedicated support",
                "White-labeling",
                "API Access",
              ]}
            />
          </div>
        )
      case "Careers":
        return (
          <div className="space-y-6">
            <h3 className="text-2xl font-bold mb-8">
              Join the remote revolution.
            </h3>
            <JobItem title="Senior Frontend Engineer" location="Remote" />
            <JobItem title="Product Marketing Manager" location="Remote" />
            <JobItem title="UX/UI Designer" location="Remote" />
          </div>
        )
      case "Integrations":
        return (
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
            <IntegrationItem
              name="Zoom"
              desc="Direct meeting sync and attendee management."
            />
            <IntegrationItem
              name="Google Meet"
              desc="Schedule via Google Calendar automatically."
            />
            <IntegrationItem
              name="Microsoft Teams"
              desc="Enterprise-level connectivity for your webinars."
            />
            <IntegrationItem
              name="Slack"
              desc="Post-event notifications and group chat sync."
            />
          </div>
        )
      default:
        return (
          <div className="prose prose-indigo max-w-none text-gray-600 leading-relaxed space-y-8">
            <h3 className="text-2xl font-black text-gray-900">
              Official {pageKey} Documentation
            </h3>
            <p>Last updated: June 1, 2024</p>
            <p>
              Welcome to the WebiStream {pageKey.toLowerCase()} page. We are
              committed to transparency and providing the best platform for our
              community.
            </p>
            <div className="space-y-4">
              <h4 className="font-bold text-gray-900">1. Overview</h4>
              <p>
                WebiStream provides high-end webinar services. This section
                covers our standards and operations regarding{" "}
                {pageKey.toLowerCase()}.
              </p>
              <h4 className="font-bold text-gray-900">2. Usage Rights</h4>
              <p>
                Users are granted a non-exclusive license to access our platform
                features according to our core principles of education and
                connectivity.
              </p>
              <h4 className="font-bold text-gray-900">3. Compliance</h4>
              <p>
                We ensure that all our {pageKey.toLowerCase()} protocols meet
                global standards, including GDPR and PCI-DSS where applicable.
              </p>
            </div>
          </div>
        )
    }
  }

  return (
    <Box minH="screen" bg="white">
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
        <Button
          variant="ghost"
          onClick={onBack}
          colorPalette="blue"
          fontWeight="bold"
          gap={2}
        >
          <ArrowLeft size={20} />
          Back
        </Button>
        <Flex align="center" gap={2}>
          <Flex bg="blue.600" p={2} borderRadius="lg">
            <Video color="white" size={20} />
          </Flex>
          <Text fontWeight="bold" fontSize="lg" letterSpacing="tight">
            WebiStream
          </Text>
        </Flex>
        <Button
          variant="ghost"
          onClick={onHostLogin}
          fontSize="sm"
          fontWeight="semibold"
          color="gray.700"
          _hover={{ color: "blue.600" }}
          px={4}
          py={2}
        >
          Host Portal
        </Button>
      </Flex>

      <Box as="main" pt={32} pb={24} px={6}>
        <Container maxW="4xl">
          <Box mb={12}>
            <Heading
              fontSize="5xl"
              fontWeight="black"
              color="gray.900"
              letterSpacing="tight"
            >
              {pageKey}
            </Heading>
            <Text color="gray.500" fontSize="lg" mt={2}>
              Detailed information about WebiStream {pageKey.toLowerCase()}{" "}
              protocols.
            </Text>
          </Box>

          {renderContent()}
        </Container>
      </Box>

      <Box
        as="footer"
        bg="gray.50"
        py={12}
        px={6}
        borderTop="1px solid"
        borderColor="gray.100"
      >
        <Container maxW="4xl" textAlign="center">
          <Text color="gray.400" fontSize="sm">
            © 2024 WebiStream Inc. All rights reserved.
          </Text>
        </Container>
      </Box>
    </Box>
  )
}

const PriceCard: React.FC<{
  tier: string
  price: string
  highlight?: boolean
  items: string[]
}> = ({ tier, price, highlight, items }) => (
  <Box
    p={8}
    borderRadius="2rem"
    border="1px solid"
    borderColor={highlight ? "blue.600" : "gray.100"}
    bg={highlight ? "blue.600" : "white"}
    color={highlight ? "white" : "inherit"}
    boxShadow={highlight ? "2xl" : "none"}
    transform={highlight ? "translateY(-8px)" : undefined}
    transition="all"
  >
    <Heading fontSize="xl" fontWeight="bold" mb={2}>
      {tier}
    </Heading>
    <Flex align="baseline" gap={1} mb={6}>
      <Text fontSize="4xl" fontWeight="black">
        {price}
      </Text>
      {price !== "Free" && price !== "Custom" && (
        <Text opacity={0.6} fontSize="sm">
          /mo
        </Text>
      )}
    </Flex>
    <Box mb={8}>
      {items.map((item, i) => (
        <Flex key={i} align="center" gap={2} fontSize="sm" mb={4}>
          <Check color={highlight ? "white" : "#4f46e5"} size={16} />
          {item}
        </Flex>
      ))}
    </Box>
    <Button
      w="full"
      py={3}
      borderRadius="2xl"
      fontWeight="bold"
      variant={highlight ? "solid" : "subtle"}
      colorPalette={highlight ? "whiteAlpha" : "blue"}
    >
      Choose {tier}
    </Button>
  </Box>
)

const JobItem: React.FC<{ title: string; location: string }> = ({
  title,
  location,
}) => (
  <Flex
    align="center"
    justify="space-between"
    p={6}
    bg="gray.50"
    borderRadius="2xl"
    cursor="pointer"
    _hover={{ bg: "blue.50" }}
    transition="all"
  >
    <Box>
      <Text
        fontWeight="bold"
        color="gray.900"
        _groupHover={{ color: "blue.600" }}
      >
        {title}
      </Text>
      <Text fontSize="sm" color="gray.500">
        {location}
      </Text>
    </Box>
    <ArrowRight size={20} color="#d1d5db" />
  </Flex>
)

const IntegrationItem: React.FC<{ name: string; desc: string }> = ({
  name,
  desc,
}) => (
  <Box
    p={6}
    border="1px solid"
    borderColor="gray.100"
    borderRadius="3xl"
    _hover={{ borderColor: "blue.100" }}
    transition="all"
  >
    <Flex align="center" justify="space-between" mb={4}>
      <Text fontWeight="bold" color="gray.900">
        {name}
      </Text>
      <ExternalLink size={16} color="#9ca3af" />
    </Flex>
    <Text fontSize="sm" color="gray.500">
      {desc}
    </Text>
  </Box>
)

export default SubPage
