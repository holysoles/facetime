-- requires facetime to already be open, in a call, with the sidebar already toggled open
-- returns list of users that were admitted to the call
set admittedUsers to {}
tell application "System Events"
	tell first window of application process "FaceTime"
		tell scroll area 2
			tell list 1 -- description "Conversation Details"
				tell list 1 -- description "X people" where X is people already in call that aren't the host
					repeat with aPossibleUser in (every group as list)
						tell aPossibleUser
							set isUser to false
							tell (first button where description is "Approve join request")
								if (exists) then
									set isUser to true
									--click it
								end if
							end tell
							if isUser then
								tell static text 1
									set admittedUsers to admittedUsers & value of it
								end tell
							end if
						end tell
					end repeat
				end tell
			end tell
		end tell
	end tell
end tell
return admittedUsers