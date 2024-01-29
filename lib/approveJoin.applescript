-- requires facetime to already be open
tell application "System Events"
	tell first window of application process "FaceTime"
		tell list 1 of list 1 of scroll area 2
			repeat with aUser in (every group as list)
				tell aUser
					repeat with checkImage in every image -- check if the element is a user row
						if description of checkImage is "contact silhouette" then
							tell (first button where description is "Approve join request")
								click
							end tell
						end if
					end repeat
				end tell
			end repeat
		end tell
	end tell
end tell