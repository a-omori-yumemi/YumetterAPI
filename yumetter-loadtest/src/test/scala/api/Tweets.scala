package api

import io.gatling.core.Predef._
import io.gatling.core.structure.ChainBuilder
import io.gatling.http.Predef._

import scala.util.Random

object Tweets {
  val get: ChainBuilder = exec(
    http("shows the timeline")
      .get("/tweets")
      .check(
        jsonPath("$[*].tweet")
          .ofType[Map[String,Any]]
          .findAll
          .saveAs("tweets")
      )
  )

  val getReplies: ChainBuilder = exec(
    http("shows replies of the tweet")
      .get("/tweets?replied_to=${tw_id}")
  )

  val post: ChainBuilder = exec(
    http("composes a new tweet")
      .post("/tweets")
      .formParamMap{ _ =>
        Map(
          "body" -> ("Hello " + Random.alphanumeric.take(20).mkString)
        )
      }
  )

  val postReply: ChainBuilder = exec(
    http("composes a new reply to a tweet")
      .post("/tweets")
      .formParamMap{ session =>
        Map(
          "body" -> "reply",
          "replied_to" -> session("tw_id").as[Int]
        )
      }
  )

  object Favourites {
    val put: ChainBuilder = exec(
      http("marks the tweet as favourite")
        .put("/tweets/${tw_id}/favorites/${usr_id}")
    )
  }
}
