package cx.sentio.botmock.config

import org.springframework.boot.web.embedded.netty.NettyReactiveWebServerFactory
import org.springframework.boot.web.reactive.server.ReactiveWebServerFactory
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration

@Configuration
open class BotMockWebFluxConfig {
    @Bean
    open fun reactiveWebServerFactory(): ReactiveWebServerFactory {
        return NettyReactiveWebServerFactory()
    }
}