package cx.sentio.botmock.config

import cx.sentio.platform.common.service.PlatformProperties
import cx.sentio.platform.common.service.ServiceProperties
import org.springframework.stereotype.Component

private const val SERVICE_NAME = "bot-mock"

@Component
data class BotMockServiceProperties(val platform: PlatformProperties) : ServiceProperties(platform, SERVICE_NAME)