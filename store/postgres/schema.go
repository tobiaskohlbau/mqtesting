//    Copyright 2017 Tobias Kohlbau
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package postgres

var migrate = []string{
	`
		create table if not exists messages (
			id            bigserial    not null primary key,
			duplicate	  boolean	   default false,
			qos			  int	       default 0,
			retained      boolean      default false,
			topic         text         not null,
			message_id    int 		   not null,
			payload       text         not null,
			created_at    timestamptz  not null
		);
	`,
}

var drop = []string{
	`drop table if exists messages cascade`,
}
