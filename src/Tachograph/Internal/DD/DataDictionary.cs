using System;

namespace Tachograph.Internal.DD
{
    /// <summary>
    /// Data Dictionary utilities for parsing tachograph data structures.
    /// Contains binary parsing logic for ~80+ data types defined in the EU regulation.
    /// </summary>
    internal static class DataDictionary
    {
        /// <summary>
        /// Unmarshal options for data dictionary types.
        /// </summary>
        internal class UnmarshalOptions
        {
            public bool PreserveRawData { get; set; } = true;
        }

        /// <summary>
        /// Marshal options for data dictionary types.
        /// </summary>
        internal class MarshalOptions
        {
            public bool UseRawData { get; set; } = true;
        }

        // TODO: Implement data dictionary parsing logic
        // - Time/date parsing (TimeReal, BCDString, Date)
        // - String encoding/decoding (code pages, IA5String)
        // - Activity records and driver identification
        // - GeoCoordinates and GNSS data
        // - Vehicle identification and calibration data
        // - ~80+ data types from EU regulation
    }
}
