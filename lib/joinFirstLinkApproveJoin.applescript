tell application "System Events"
	tell application process "FaceTime"
		tell first window -- default window
			tell scroll area 2 -- unlabelled
				tell (first list where description is "Recent Calls")
					tell (first list where description is "Upcoming")
						tell first group -- most recently created call/link
							repeat with anAction in (actions as list)
								if (description of anAction) is "FaceTime Video" then
									set desiredAction to anAction
									exit repeat
								end if
							end repeat
							perform desiredAction
						end tell
					end tell
				end tell
			end tell
			delay 2 -- feels unnecessary be necessary but makes joining more reliable. Seems like the camera needs to initialize first
			tell button "Join"
				repeat until (exists)
					delay 0.1
				end repeat
				click
			end tell
			tell list 1 of list 1 of scroll area 2 -- poll since this needs to time to populate
				repeat until (exists)
					delay 0.1
				end repeat
				repeat with aPossibleUser in (every group as list)
					tell aPossibleUser
						repeat with checkImage in every image -- validate the element is a user row
							if description of checkImage is "contact silhouette" then
								tell (first button where description is "Approve join request") --this takes a beat to get displayed
									repeat until (exists)
										delay 0.1
									end repeat
									click
								end tell
							end if
						end repeat
					end tell
				end repeat
			end tell
		end tell
	end tell
end tell