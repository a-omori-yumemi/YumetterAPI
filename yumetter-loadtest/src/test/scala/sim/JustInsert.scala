package sim

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import scala.concurrent.duration._
import scala.util.Random
import scala.math.exp
import java.util.Calendar


class JustInsert extends Simulation {

  val httpProtocol = http
    .baseUrl(sys.env("YUMETTER_API") + "/v1") // Here is the root for all relative URLs
    .acceptHeader("application/json;q=0.9,*/*;q=0.8") // Here are the common headers
    .acceptEncodingHeader("gzip, deflate")
    .acceptLanguageHeader("en-US,en;q=0.5")
    .userAgentHeader("Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:16.0) Gecko/20100101 Firefox/16.0")

  val scn = scenario("Insert scenario") // A scenario is a chain of requests and pauses
    .exec(RegisterAndLogin.reg)
    .repeat(2000)(
      exec(Tweet.tweet)
      .exec(WatchTL.watchFavReply)
    )

  setUp(scn.inject(rampUsers(10).during(30.seconds)).protocols(httpProtocol))
}
