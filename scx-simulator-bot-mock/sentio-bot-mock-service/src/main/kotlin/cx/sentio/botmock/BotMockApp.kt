package cx.sentio.botmock

import org.springframework.boot.WebApplicationType
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.builder.SpringApplicationBuilder
import org.springframework.web.reactive.config.EnableWebFlux

@SpringBootApplication
@EnableWebFlux
open class BotMockApp {
    companion object {
        @JvmStatic
        fun main(args: Array<String>) {
            SpringApplicationBuilder()
                .sources(BotMockApp::class.java)
                .web(WebApplicationType.REACTIVE)
                .run(*args)
        }
    }
}